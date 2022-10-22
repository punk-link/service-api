package platforms

import (
	"encoding/json"
	platformData "main/data/platforms"
	"main/models/platforms"
	"main/services/artists"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/punk-link/logger"
	platformContracts "github.com/punk-link/platform-contracts"
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

func (t *StreamingPlatformService) ProcessPlatforeUrlResults() {
	jetStreamContext, err := t.natsConnection.JetStream()
	err = t.createReducerJetStreamIfNotExist(err, jetStreamContext)
	subscription, err := t.getSubscription(err, jetStreamContext)
	t.consumeUrlResults(err, subscription)
}

func (t *StreamingPlatformService) PublishPlatforeUrlRequests() {
	jetStreamContext, err := t.natsConnection.JetStream()
	err = t.createPlatformJetStreamIfNotExist(err, jetStreamContext)
	err = t.publishPlatforeUrlRequests(err, jetStreamContext)
	if err != nil {
		t.logger.LogError(err, err.Error())
	}
}

func (t *StreamingPlatformService) consumeUrlResults(err error, subscription *nats.Subscription) error {
	if err != nil {
		return err
	}

	for {
		messages, _ := subscription.Fetch(ITERATION_STEP)
		urlResults := make([]platformContracts.UrlResultContainer, len(messages))
		for i, message := range messages {
			message.Ack()

			var urlResult platformContracts.UrlResultContainer
			_ = json.Unmarshal(message.Data, &urlResult)

			urlResults[i] = urlResult
		}

		t.resync(urlResults)
	}
}

func (t *StreamingPlatformService) createPlatformJetStreamIfNotExist(err error, jetStreamContext nats.JetStreamContext) error {
	if err != nil {
		return err
	}

	stream, _ := jetStreamContext.StreamInfo(platformContracts.PLATFORM_URL_REQUESTS_STREAM_NAME)
	if stream == nil {
		t.logger.LogInfo("Creating Nats stream %s and subjects %s", platformContracts.PLATFORM_URL_REQUESTS_STREAM_NAME, platformContracts.PLATFORM_URL_REQUESTS_STREAM_SUBJECTS)
		_, err = jetStreamContext.AddStream(platformContracts.DefaultPlatformServiceConfig)
	}

	return err
}

func (t *StreamingPlatformService) createReducerJetStreamIfNotExist(err error, jetStreamContext nats.JetStreamContext) error {
	if err != nil {
		return err
	}

	stream, _ := jetStreamContext.StreamInfo(platformContracts.PLATFORM_URL_RESPONSE_STREAM_NAME)
	if stream == nil {
		t.logger.LogInfo("Creating Nats stream %s and subjects %s", platformContracts.PLATFORM_URL_RESPONSE_STREAM_NAME, platformContracts.PLATFORM_URL_RESPONSE_STREAM_SUBJECT)
		_, err = jetStreamContext.AddStream(platformContracts.DefaultReducerConfig)
	}

	return err
}

func (t *StreamingPlatformService) getExistedUrls(logger logger.Logger, urlResults []platformContracts.UrlResultContainer) (map[int]platformData.PlatformReleaseUrl, error) {
	ids := make([]int, len(urlResults))
	for i, result := range urlResults {
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

func (t *StreamingPlatformService) getPlatformReleaseUrls(err error, existedUrls map[int]platformData.PlatformReleaseUrl, upcResults []platformContracts.UrlResultContainer, timestamp time.Time) ([]platformData.PlatformReleaseUrl, error) {
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

func (t *StreamingPlatformService) getSubscription(err error, jetStreamContext nats.JetStreamContext) (*nats.Subscription, error) {
	if err != nil {
		return nil, err
	}

	return jetStreamContext.PullSubscribe(platformContracts.PLATFORM_URL_RESPONSE_STREAM_SUBJECT, platformContracts.PLATFORM_URL_RESPONSE_CONSUMER_NAME)
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
		for _, platform := range platformContracts.AvailablePlatforms {
			subjectName := platformContracts.GetRequestStreamSubject(platform)
			for _, container := range upcContainers {
				json, _ := json.Marshal(container)
				jetStreamContext.Publish(subjectName, json)
			}
		}

		skip += ITERATION_STEP
	}

	return err
}

func (t *StreamingPlatformService) resync(urlResults []platformContracts.UrlResultContainer) {
	timestamp := time.Now().UTC()

	existedUrls, err := t.getExistedUrls(t.logger, urlResults)
	platformReleaseUrls, err := t.getPlatformReleaseUrls(err, existedUrls, urlResults, timestamp)
	newUrls, changedUrls := distinctNewUrlsFromChanged(platformReleaseUrls)

	err = createDbPlatformReleaseUrlsInBatches(t.logger, err, newUrls)
	err = updateDbPlatformReleaseUrlsInBatches(t.logger, err, changedUrls)
	t.markReleasesAsUpdated(err, platformReleaseUrls, timestamp)
}

func buildChangedPlatformReleaseUrl(existedUrl platformData.PlatformReleaseUrl, urlResult platformContracts.UrlResultContainer, timestamp time.Time) platformData.PlatformReleaseUrl {
	return platformData.PlatformReleaseUrl{
		Id:      existedUrl.Id,
		Updated: timestamp,
		Url:     urlResult.Url,
	}
}

func buildNewPlatformReleaseUrl(urlResult platformContracts.UrlResultContainer, timestamp time.Time) platformData.PlatformReleaseUrl {
	return platformData.PlatformReleaseUrl{
		Created:      timestamp,
		PlatformName: urlResult.PlatformName,
		ReleaseId:    urlResult.Id,
		Updated:      timestamp,
		Url:          urlResult.Url,
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
