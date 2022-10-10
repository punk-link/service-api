package deezer

import (
	platforms "main/models/platforms"
	deezerModels "main/models/platforms/deezer"
	platformEnums "main/models/platforms/enums"
	"main/services/common"
	platformServices "main/services/platforms/base"

	"github.com/samber/do"
)

type DeezerService struct {
	logger *common.Logger
}

func ConstructDeezerService(injector *do.Injector) (*DeezerService, error) {
	logger := do.MustInvoke[*common.Logger](injector)

	return &DeezerService{
		logger: logger,
	}, nil
}

func (t *DeezerService) GetPlatformName() string {
	return platformEnums.Deezer
}

func (t *DeezerService) GetReleaseUrlsByUpc(upcContainers []platforms.UpcContainer) []platforms.UrlResultContainer {
	syncedUrls := platformServices.GetSyncedUrls("album/upc:%s", upcContainers)
	upcMap := platformServices.GetUpcMap(upcContainers)

	syncedUpcResults := makeBatchRequestWithSync[deezerModels.UpcResponse](t.logger, "GET", syncedUrls)

	results := make([]platforms.UrlResultContainer, len(syncedUpcResults))
	for i, syncedResult := range syncedUpcResults {
		id := upcMap[syncedResult.Sync]
		results[i] = platformServices.BuildUrlResultContainer(id, t.GetPlatformName(), syncedResult.Sync, syncedResult.Result.Url)
	}

	return results
}
