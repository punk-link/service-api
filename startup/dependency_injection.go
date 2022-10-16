package startup

import (
	"fmt"
	apiControllers "main/controllers/api"
	staticControllers "main/controllers/static"
	platformConstants "main/models/platforms/constants"
	platformEnums "main/models/platforms/enums"
	"main/models/platforms/spotify/accessToken"
	artistServices "main/services/artists"
	artistStaticServices "main/services/artists/static"
	"main/services/cache"
	"main/services/common"
	loggerServices "main/services/common/logger"
	labelServices "main/services/labels"
	platformServices "main/services/platforms"
	deezerServices "main/services/platforms/deezer"
	spotifyServices "main/services/platforms/spotify"

	consulClient "github.com/punk-link/consul-client"

	"github.com/punk-link/logger"
	"github.com/samber/do"
)

func buildDependencies(logger *logger.Logger, consul *consulClient.ConsulClient) *do.Injector {
	container := do.New()

	spotifySettingsValue, err := consul.GetOrSet("SpotifySettings", 0)
	if err != nil {
		logger.LogFatal(err, "Can't obtain Spotify settings from Consul: %s", err.Error())
	}

	spotifySettings := spotifySettingsValue.(map[string]interface{})
	do.ProvideValue(container, &accessToken.SpotifyClientConfig{
		ClientId:     spotifySettings["ClientId"].(string),
		ClientSecret: spotifySettings["ClientSecret"].(string),
	})

	do.Provide(container, loggerServices.New)
	do.Provide(container, common.ConstructHashCoder)
	do.Provide(container, cache.ConstructMemoryCacheService)

	do.Provide(container, labelServices.ConstructLabelService)
	do.Provide(container, labelServices.ConstructManagerService)

	do.Provide(container, deezerServices.ConstructDeezerService)
	do.ProvideNamed(container, buildPlatformServiceName(platformEnums.Deezer), deezerServices.ConstructDeezerServiceAsPlatformer)

	do.Provide(container, spotifyServices.ConstructSpotifyService)
	do.ProvideNamed(container, buildPlatformServiceName(platformEnums.Spotify), spotifyServices.ConstructSpotifyServiceAsPlatformer)

	do.Provide(container, artistServices.ConstructReleaseService)
	do.Provide(container, artistStaticServices.ConstructStaticReleaseService)
	do.Provide(container, artistServices.ConstructArtistService)
	do.Provide(container, artistStaticServices.ConstructStaticArtistService)

	do.Provide(container, platformServices.ConstructStreamingPlatformService)

	do.Provide(container, apiControllers.ConstructArtistController)
	do.Provide(container, apiControllers.ConstructHashController)
	do.Provide(container, apiControllers.ConstructLabelController)
	do.Provide(container, apiControllers.ConstructManagerController)
	do.Provide(container, apiControllers.ConstructStreamingPlatformController)
	do.Provide(container, apiControllers.ConstructReleaseController)
	do.Provide(container, apiControllers.ConstructStatusController)

	do.Provide(container, staticControllers.ConstructStaticArtistController)
	do.Provide(container, staticControllers.ConstructStaticReleaseController)

	return container
}

func buildPlatformServiceName(serviceName string) string {
	return fmt.Sprintf("%s%s", serviceName, platformConstants.PLATFORM_SERVICE_TOKEN)
}
