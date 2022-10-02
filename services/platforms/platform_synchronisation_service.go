package platforms

import (
	"fmt"
	platformData "main/data/platforms"
	"main/models/platforms"
	basePlatforms "main/models/platforms/base"
	platformConstants "main/models/platforms/constants"
	"main/services/artists"
	"main/services/common"
	"main/services/platforms/base"
	"sync"
	"time"

	"github.com/samber/do"
)

type PlatformSynchronisationService struct {
	injector       *do.Injector
	logger         *common.Logger
	releaseService *artists.ReleaseService
}

func ConstructPlatformSynchronisationService(injector *do.Injector) (*PlatformSynchronisationService, error) {
	logger := do.MustInvoke[*common.Logger](injector)
	releaseService := do.MustInvoke[*artists.ReleaseService](injector)

	return &PlatformSynchronisationService{
		injector:       injector.Clone(),
		logger:         logger,
		releaseService: releaseService,
	}, nil
}

func (t *PlatformSynchronisationService) Sync() {
	now := time.Now().UTC()

	urls := t.getPlatformReleaseUrlsToSync(now)
	t.resync(urls, now)
}

func (t *PlatformSynchronisationService) getExistedUrls(logger *common.Logger, upcContainers []platforms.UpcContainer, upcResults []platforms.UrlResultContainer) (map[int]platformData.PlatformReleaseUrl, error) {
	ids := make([]int, len(upcContainers))
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

func (t *PlatformSynchronisationService) getPlatformerContainers() []basePlatforms.PlatformerContainer {
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

func (t *PlatformSynchronisationService) getPlatformReleaseUrls(err error, existedUrls map[int]platformData.PlatformReleaseUrl, upcResults []platforms.UrlResultContainer, timestamp time.Time) ([]platformData.PlatformReleaseUrl, error) {
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

func (t *PlatformSynchronisationService) getPlatformReleaseUrlsToSync(timestamp time.Time) []platformData.PlatformReleaseUrl {
	releaseCount := t.releaseService.GetCount()
	updateTreshold := time.Now().UTC().Add(-UPDATE_TRESHOLD_MINUTES)

	platformerContainers := t.getPlatformerContainers()

	var wg sync.WaitGroup
	chanResults := make(chan []platformData.PlatformReleaseUrl)

	skip := 0
	for i := 0; i < releaseCount; i = i + ITERATION_STEP {
		upcContainers := t.releaseService.GetUpcContainersToUpdate(ITERATION_STEP, skip, updateTreshold)

		wg.Add(1)
		go t.getUrlsToResync(&wg, chanResults, platformerContainers, upcContainers, timestamp)

		skip += ITERATION_STEP
	}

	go func() {
		wg.Wait()
		close(chanResults)
	}()

	urls := make([]platformData.PlatformReleaseUrl, 0)
	for result := range chanResults {
		urls = append(urls, result...)
	}

	return urls
}

func (t *PlatformSynchronisationService) getUpcResultsFromPlatformers(platformerContainers []basePlatforms.PlatformerContainer, upcContainers []platforms.UpcContainer) []platforms.UrlResultContainer {
	upcResults := make([]platforms.UrlResultContainer, 0)
	for _, container := range platformerContainers {
		platformer := container.Instance
		partialUpcResults := platformer.GetReleaseUrlsByUpc(upcContainers)

		upcResults = append(upcResults, partialUpcResults...)
	}

	return upcResults
}

func (t *PlatformSynchronisationService) getUrlsToResync(wg *sync.WaitGroup, results chan<- []platformData.PlatformReleaseUrl, platformerContainers []basePlatforms.PlatformerContainer, upcContainers []platforms.UpcContainer, timestamp time.Time) {
	defer wg.Done()

	upcResults := t.getUpcResultsFromPlatformers(platformerContainers, upcContainers)
	existedUrls, err := t.getExistedUrls(t.logger, upcContainers, upcResults)
	platformReleaseUrls, err := t.getPlatformReleaseUrls(err, existedUrls, upcResults, timestamp)
	if err != nil {
		return
	}

	results <- platformReleaseUrls
}

func (t *PlatformSynchronisationService) markReleasesAsUpdated(err error, platformReleaseUrls []platformData.PlatformReleaseUrl, timestamp time.Time) error {
	if err != nil {
		return err
	}

	releaseIds := make([]int, len(platformReleaseUrls))
	for i, url := range platformReleaseUrls {
		releaseIds[i] = url.ReleaseId
	}

	return t.releaseService.MarkAsUpdated(releaseIds, timestamp)
}

func (t *PlatformSynchronisationService) resync(platformReleaseUrls []platformData.PlatformReleaseUrl, timestamp time.Time) {
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
const UPDATE_TRESHOLD_MINUTES = time.Minute * 0
