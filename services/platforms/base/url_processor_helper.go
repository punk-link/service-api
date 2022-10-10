package base

import (
	"fmt"
	commonModels "main/models/common"
	platforms "main/models/platforms"
)

func BuildUrlResultContainer(id int, platformName string, upc string, url string) platforms.UrlResultContainer {
	return platforms.UrlResultContainer{
		Id:           id,
		PlatformName: platformName,
		Upc:          upc,
		Url:          url,
	}
}

func GetSyncedUrls(format string, upcContainers []platforms.UpcContainer) []commonModels.SyncedUrl {
	results := make([]commonModels.SyncedUrl, len(upcContainers))
	for i, container := range upcContainers {
		results[i] = commonModels.SyncedUrl{
			Sync: container.Upc,
			Url:  fmt.Sprintf(format, container.Upc),
		}
	}

	return results
}

func GetUpcMap(upcContainers []platforms.UpcContainer) map[string]int {
	results := make(map[string]int, len(upcContainers))
	for _, container := range upcContainers {
		results[container.Upc] = container.Id
	}

	return results
}
