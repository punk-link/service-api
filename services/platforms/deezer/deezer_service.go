package deezer

import (
	"fmt"
	platforms "main/models/platforms"
	deezerModels "main/models/platforms/deezer"
	platformEnums "main/models/platforms/enums"
	platformServices "main/services/platforms/base"
	"time"

	"github.com/punk-link/logger"

	"github.com/samber/do"
)

type DeezerService struct {
	logger *logger.Logger
}

func ConstructDeezerService(injector *do.Injector) (*DeezerService, error) {
	logger := do.MustInvoke[*logger.Logger](injector)

	return &DeezerService{
		logger: logger,
	}, nil
}

func ConstructDeezerServiceAsPlatformer(injector *do.Injector) (platformServices.Platformer, error) {
	logger := do.MustInvoke[*logger.Logger](injector)

	return platformServices.Platformer(&DeezerService{
		logger: logger,
	}), nil
}

func (t *DeezerService) GetPlatformName() string {
	return platformEnums.Deezer
}

func (t *DeezerService) GetReleaseUrlsByUpc(upcContainers []platforms.UpcContainer) []platforms.UrlResultContainer {
	results := make([]platforms.UrlResultContainer, 0)
	for _, container := range upcContainers {
		var response deezerModels.UpcResponse
		err := makeRequest(t.logger, "GET", fmt.Sprintf("album/upc:%s", container.Upc), &response)
		if err != nil {
			continue
		}

		if response.Error.Code != 0 {
			continue
		}

		results = append(results, platformServices.BuildUrlResultContainer(container.Id, t.GetPlatformName(), container.Upc, response.Url))

		time.Sleep(REQUEST_TIMEOUT_DURATION_MILLISECONDS)
	}

	return results
}

const REQUEST_TIMEOUT_DURATION_MILLISECONDS = time.Millisecond * 100
