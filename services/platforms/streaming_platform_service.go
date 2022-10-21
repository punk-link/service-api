package platforms

import (
	"encoding/json"
	"fmt"
	platformData "main/data/platforms"
	"main/models/platforms"
	basePlatforms "main/models/platforms/base"
	platformConstants "main/models/platforms/constants"
	"main/services/artists"
	"main/services/platforms/base"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type StreamingPlatformService struct {
	injector       *do.Injector
	logger         logger.Logger
	natsConnection *nats.Conn
	releaseService *artists.ReleaseService
}

func ConstructStreamingPlatformService(injector *do.Injector) (*StreamingPlatformService, error) {
	logger := do.MustInvoke[logger.Logger](injector)
	natsConnection := do.MustInvoke[*nats.Conn](injector)
	releaseService := do.MustInvoke[*artists.ReleaseService](injector)

	return &StreamingPlatformService{
		injector:       injector.Clone(),
		logger:         logger,
		natsConnection: natsConnection,
		releaseService: releaseService,
	}, nil
}

func (t *StreamingPlatformService) Get(releaseId int) ([]platforms.PlatformReleaseUrl, error) {
	urls, err := getDbPlatformReleaseUrlsByReleaseId(t.logger, nil, releaseId)
	if err != nil {
		return make([]platforms.PlatformReleaseUrl, 0), err
	}

	results := make([]platforms.PlatformReleaseUrl, len(urls))
	for i, url := range urls {
		results[i] = platforms.PlatformReleaseUrl{
			Id:           url.Id,
			PlatformName: url.PlatformName,
			ReleaseId:    url.ReleaseId,
			Url:          url.Url,
		}
	}

	return results, err
}

func (t *StreamingPlatformService) PublishPlatforeUrlRequests() {
	jetStreamContext, err := t.natsConnection.JetStream()
	err = t.createJstStreamIfNotExist(err, jetStreamContext)
	err = t.publishPlatforeUrlRequests(err, jetStreamContext)
	if err != nil {
		t.logger.LogError(err, err.Error())
	}
}

func (t *StreamingPlatformService) createJstStreamIfNotExist(err error, jetStreamContext nats.JetStreamContext) error {
	if err != nil {
		return err
	}

	stream, _ := jetStreamContext.StreamInfo(PLATFORM_URL_REQUESTS_STREAM_NAME)
	if stream == nil {
		t.logger.LogInfo("Creating Nats stream %s and subjects %s", PLATFORM_URL_REQUESTS_STREAM_NAME, PLATFORM_URL_REQUESTS_STREAM_SUBJECTS)
		_, err = jetStreamContext.AddStream(&nats.StreamConfig{
			Name:      PLATFORM_URL_REQUESTS_STREAM_NAME,
			MaxAge:    time.Hour,
			Retention: nats.WorkQueuePolicy,
			Storage:   nats.MemoryStorage,
			Subjects:  []string{PLATFORM_URL_REQUESTS_STREAM_SUBJECTS},
		})
	}

	return err
}

func (t *StreamingPlatformService) publishPlatforeUrlRequests(err error, jetStreamContext nats.JetStreamContext) error {
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	releaseCount := t.releaseService.GetCount()
	updateTreshold := now.Add(-UPDATE_TRESHOLD_INTERVAL)

	skip := 0
	for i := 0; i < releaseCount; i = i + ITERATION_STEP {
		upcContainers := t.releaseService.GetUpcContainersToUpdate(ITERATION_STEP, skip, updateTreshold)
		for _, platform := range platforms.AvailablePlatforms {
			subjectName := fmt.Sprintf("%s.%s", PLATFORM_URL_REQUESTS_STREAM_NAME, platform)
			for _, container := range upcContainers {
				json, _ := json.Marshal(container)
				jetStreamContext.Publish(subjectName, json)
			}
		}

		skip += ITERATION_STEP
	}

	return err
}

const PLATFORM_URL_REQUESTS_STREAM_NAME = "PLATFORM-URL-REQUESTS"
const PLATFORM_URL_REQUESTS_STREAM_SUBJECTS = "PLATFORM-URL-REQUESTS.*"

func (t *StreamingPlatformService) getExistedUrls(logger logger.Logger, upcResults []platforms.UrlResultContainer) (map[int]platformData.PlatformReleaseUrl, error) {
	ids := make([]int, len(upcResults))
	for i, result := range upcResults {
		ids[i] = result.Id
	}

	existedUrls, err := getDbPlatformReleaseUrlsByReleaseIds(logger, nil, ids)
	if err != nil {
		return make(map[int]platformData.PlatformReleaseUrl, 0), err
	}

	existedUrlsMap := make(map[int]platformData.PlatformReleaseUrl, len(existedUrls))
	for _, url := range existedUrls {
		existedUrlsMap[url.ReleaseId] = url
	}

	return existedUrlsMap, err
}

func (t *StreamingPlatformService) getPlatformerContainers() []basePlatforms.PlatformerContainer {
	platformerContainers := make([]basePlatforms.PlatformerContainer, len(platforms.AvailablePlatforms))
	for i, platformName := range platforms.AvailablePlatforms {
		fullPlatformName := fmt.Sprintf("%s%s", platformName, platformConstants.PLATFORM_SERVICE_TOKEN)
		platformer := do.MustInvokeNamed[base.Platformer](t.injector, fullPlatformName)

		platformerContainers[i] = basePlatforms.PlatformerContainer{
			Instance:        platformer,
			FullServiceName: fullPlatformName,
		}
	}

	return platformerContainers
}

func (t *StreamingPlatformService) getPlatformReleaseUrls(err error, existedUrls map[int]platformData.PlatformReleaseUrl, upcResults []platforms.UrlResultContainer, timestamp time.Time) ([]platformData.PlatformReleaseUrl, error) {
	if err != nil {
		return make([]platformData.PlatformReleaseUrl, 0), err
	}

	platformReleaseUrls := make([]platformData.PlatformReleaseUrl, 0)
	for _, upcResult := range upcResults {
		if existedUrl, isExist := existedUrls[upcResult.Id]; isExist {
			if existedUrl.Url != upcResult.Url {
				platformReleaseUrls = append(platformReleaseUrls, buildChangedPlatformReleaseUrl(existedUrl, upcResult, timestamp))
			}

			continue
		}

		platformReleaseUrls = append(platformReleaseUrls, buildNewPlatformReleaseUrl(upcResult, timestamp))
	}

	return platformReleaseUrls, err
}

func (t *StreamingPlatformService) getUpcResultsFromPlatformers(platformerContainers []basePlatforms.PlatformerContainer, upcContainers []platforms.UpcContainer) []platforms.UrlResultContainer {
	upcResults := make([]platforms.UrlResultContainer, 0)
	for _, container := range platformerContainers {
		platformer := container.Instance
		partialUpcResults := platformer.GetReleaseUrlsByUpc(upcContainers)

		upcResults = append(upcResults, partialUpcResults...)
	}

	return upcResults
}

func (t *StreamingPlatformService) getUrlsToResync(wg *sync.WaitGroup, results chan<- []platformData.PlatformReleaseUrl, platformerContainers []basePlatforms.PlatformerContainer, upcContainers []platforms.UpcContainer, timestamp time.Time) {
	defer wg.Done()

	upcResults := t.getUpcResultsFromPlatformers(platformerContainers, upcContainers)
	existedUrls, err := t.getExistedUrls(t.logger, upcResults)
	platformReleaseUrls, err := t.getPlatformReleaseUrls(err, existedUrls, upcResults, timestamp)
	if err != nil {
		return
	}

	results <- platformReleaseUrls
}

func (t *StreamingPlatformService) markReleasesAsUpdated(err error, platformReleaseUrls []platformData.PlatformReleaseUrl, timestamp time.Time) error {
	if err != nil {
		return err
	}

	releaseIds := make([]int, len(platformReleaseUrls))
	for i, url := range platformReleaseUrls {
		releaseIds[i] = url.ReleaseId
	}

	return t.releaseService.MarkAsUpdated(releaseIds, timestamp)
}

func (t *StreamingPlatformService) resync(platformReleaseUrls []platformData.PlatformReleaseUrl, timestamp time.Time) {
	newUrls, changedUrls := distinctNewUrlsFromChanged(platformReleaseUrls)

	err := createDbPlatformReleaseUrlsInBatches(t.logger, nil, newUrls)
	err = updateDbPlatformReleaseUrlsInBatches(t.logger, err, changedUrls)
	t.markReleasesAsUpdated(err, platformReleaseUrls, timestamp)
}

func buildChangedPlatformReleaseUrl(existedUrl platformData.PlatformReleaseUrl, upcResult platforms.UrlResultContainer, timestamp time.Time) platformData.PlatformReleaseUrl {
	return platformData.PlatformReleaseUrl{
		Id:      existedUrl.Id,
		Updated: timestamp,
		Url:     upcResult.Url,
	}
}

func buildNewPlatformReleaseUrl(upcResult platforms.UrlResultContainer, timestamp time.Time) platformData.PlatformReleaseUrl {
	return platformData.PlatformReleaseUrl{
		Created:      timestamp,
		PlatformName: upcResult.PlatformName,
		ReleaseId:    upcResult.Id,
		Updated:      timestamp,
		Url:          upcResult.Url,
	}
}

func distinctNewUrlsFromChanged(platformReleaseUrls []platformData.PlatformReleaseUrl) ([]platformData.PlatformReleaseUrl, []platformData.PlatformReleaseUrl) {
	changedUrls := make([]platformData.PlatformReleaseUrl, 0)
	newUrls := make([]platformData.PlatformReleaseUrl, 0)
	for _, result := range platformReleaseUrls {
		if result.Id == 0 {
			newUrls = append(newUrls, result)
		} else {
			changedUrls = append(changedUrls, result)
		}
	}

	return newUrls, changedUrls
}

const ITERATION_STEP = 40
const UPDATE_TRESHOLD_INTERVAL = time.Hour
