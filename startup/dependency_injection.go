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

	"github.com/nats-io/nats.go"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

func buildDependencies(logger *logger.Logger, consul *consulClient.ConsulClient) *do.Injector {
	injector := do.New()

	spotifySettingsValue, err := consul.GetOrSet("SpotifySettings", 0)
	if err != nil {
		logger.LogFatal(err, "Can't obtain Spotify settings from Consul: %s", err.Error())
	}

	spotifySettings := spotifySettingsValue.(map[string]interface{})
	do.ProvideValue(injector, &accessToken.SpotifyClientConfig{
		ClientId:     spotifySettings["ClientId"].(string),
		ClientSecret: spotifySettings["ClientSecret"].(string),
	})

	natsConnection, err := nats.Connect(nats.DefaultOptions.Url)
	if err != nil {
		logger.LogFatal(err, "Nats connection error: %s", err.Error())
	}

	do.ProvideValue(injector, natsConnection)

	do.Provide(injector, loggerServices.New)
	do.Provide(injector, common.ConstructHashCoder)
	do.Provide(injector, cache.ConstructMemoryCacheService)

	do.Provide(injector, labelServices.ConstructLabelService)
	do.Provide(injector, labelServices.ConstructManagerService)

	do.Provide(injector, deezerServices.ConstructDeezerService)
	do.ProvideNamed(injector, buildPlatformServiceName(platformEnums.Deezer), deezerServices.ConstructDeezerServiceAsPlatformer)

	do.Provide(injector, spotifyServices.ConstructSpotifyService)
	do.ProvideNamed(injector, buildPlatformServiceName(platformEnums.Spotify), spotifyServices.ConstructSpotifyServiceAsPlatformer)

	do.Provide(injector, artistServices.ConstructReleaseService)
	do.Provide(injector, artistStaticServices.ConstructStaticReleaseService)
	do.Provide(injector, artistServices.ConstructArtistService)
	do.Provide(injector, artistStaticServices.ConstructStaticArtistService)

	do.Provide(injector, platformServices.ConstructStreamingPlatformService)

	do.Provide(injector, apiControllers.ConstructArtistController)
	do.Provide(injector, apiControllers.ConstructHashController)
	do.Provide(injector, apiControllers.ConstructLabelController)
	do.Provide(injector, apiControllers.ConstructManagerController)
	do.Provide(injector, apiControllers.ConstructStreamingPlatformController)
	do.Provide(injector, apiControllers.ConstructReleaseController)
	do.Provide(injector, apiControllers.ConstructStatusController)

	do.Provide(injector, staticControllers.ConstructStaticArtistController)
	do.Provide(injector, staticControllers.ConstructStaticReleaseController)

	return injector
}

func buildPlatformServiceName(serviceName string) string {
	return fmt.Sprintf("%s%s", serviceName, platformConstants.PLATFORM_SERVICE_TOKEN)
}
