package platforms

import (
	"fmt"
	"main/models/platforms"
	basePlatforms "main/models/platforms/base"
	platformConstants "main/models/platforms/constants"
	"main/services/artists"
	"main/services/platforms/base"
	"time"

	"github.com/samber/do"
)

type PlatformSynchronisationService struct {
	injector       *do.Injector
	releaseService *artists.ReleaseService
}

func ConstructPlatformSynchronisationService(injector *do.Injector) (*PlatformSynchronisationService, error) {
	releaseService := do.MustInvoke[*artists.ReleaseService](injector)

	return &PlatformSynchronisationService{
		injector:       injector.Clone(),
		releaseService: releaseService,
	}, nil
}

func (t *PlatformSynchronisationService) Sync() {
	platformerContainers := t.getPlatformerContainers()

	updateTreshold := time.Now().UTC().Add(-UPDATE_TRESHOLD_MINUTES)
	releaseCount := t.releaseService.GetCount()

	skip := 0
	for i := 0; i < releaseCount; i = i + ITERATION_STEP {
		upcContainers := t.releaseService.GetUpcContainersToUpdate(ITERATION_STEP, skip, updateTreshold)

		for _, container := range platformerContainers {
			platformer := container.Instance
			platformer.GetReleaseUrlsByUpc(upcContainers)

		}

		skip += ITERATION_STEP
	}

	t.disposePlatformers(platformerContainers)
}

func (t *PlatformSynchronisationService) disposePlatformers(platformerContainers []basePlatforms.PlatformerContainer) {
	for _, container := range platformerContainers {
		do.ShutdownNamed(t.injector, container.FullServiceName)
	}
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

const ITERATION_STEP = 40
const UPDATE_TRESHOLD_MINUTES = time.Minute * 15
