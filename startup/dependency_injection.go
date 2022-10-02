package startup

import (
	"fmt"
	apiControllers "main/controllers/api"
	mvcControllers "main/controllers/mvc"
	platformConstants "main/models/platforms/constants"
	platformEnums "main/models/platforms/enums"
	artistServices "main/services/artists"
	artistMvcServices "main/services/artists/mvc"
	"main/services/cache"
	"main/services/common"
	labelServices "main/services/labels"
	platformServices "main/services/platforms"
	spotifyServices "main/services/spotify"

	"github.com/samber/do"
)

func buildDependencies() *do.Injector {
	container := do.New()

	do.Provide(container, common.ConstructLogger)
	do.Provide(container, common.ConstructHashCoder)
	do.Provide(container, cache.ConstructMemoryCacheService)

	do.Provide(container, labelServices.ConstructLabelService)
	do.Provide(container, labelServices.ConstructManagerService)

	do.Provide(container, spotifyServices.ConstructSpotifyService)
	do.ProvideNamed(container, buildPlatformServiceName(platformEnums.Spotify), spotifyServices.ConstructSpotifyServiceAsPlatformer)

	do.Provide(container, artistServices.ConstructReleaseService)
	do.Provide(container, artistMvcServices.ConstructMvcReleaseService)
	do.Provide(container, artistServices.ConstructArtistService)
	do.Provide(container, artistMvcServices.ConstructMvcArtistService)

	do.Provide(container, platformServices.ConstructPlatformSynchronisationService)

	do.Provide(container, apiControllers.ConstructArtistController)
	do.Provide(container, apiControllers.ConstructHashController)
	do.Provide(container, apiControllers.ConstructLabelController)
	do.Provide(container, apiControllers.ConstructManagerController)
	do.Provide(container, apiControllers.ConstructPlatformSynchronisationController)
	do.Provide(container, apiControllers.ConstructReleaseController)
	do.Provide(container, apiControllers.ConstructStatusController)

	do.Provide(container, mvcControllers.ConstructMvcArtistController)
	do.Provide(container, mvcControllers.ConstructMvcReleaseController)

	return container
}

func buildPlatformServiceName(serviceName string) string {
	return fmt.Sprintf("%s%s", serviceName, platformConstants.PLATFORM_SERVICE_TOKEN)
}
