package platforms

import (
	"fmt"
	"main/models/platforms"
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
	updateTreshold := time.Now().UTC().Add(-UPDATE_TRESHOLD_MINUTES)

	releaseCount := t.releaseService.GetCount()

	skip := 0
	for i := 0; i < releaseCount; i = i + ITERATION_STEP {
		upcContainers := t.releaseService.GetUpcContainersToUpdate(ITERATION_STEP, skip, updateTreshold)

		for _, platformName := range platforms.AvailablePlatforms {
			platform := do.MustInvokeNamed[base.Platformer](t.injector, fmt.Sprintf("%s%s", platformName, platformConstants.PLATFORM_SERVICE_TOKEN))
			platform.GetReleaseUrlsByUpc(upcContainers)
		}

		skip += ITERATION_STEP
	}

	// for _, platform := range platforms.AvailablePlatforms {

	// }
}

const ITERATION_STEP = 40
const UPDATE_TRESHOLD_MINUTES = time.Minute * 15
