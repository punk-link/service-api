package startup

import (
	apiControllers "main/controllers/api"
	staticControllers "main/controllers/static"
	"main/models/platforms/spotify/accessToken"
	artistServices "main/services/artists"
	artistStaticServices "main/services/artists/static"
	"main/services/cache"
	"main/services/common"
	loggerServices "main/services/common/logger"
	labelServices "main/services/labels"
	platformServices "main/services/platforms"
	spotifyServices "main/services/platforms/spotify"

	consulClient "github.com/punk-link/consul-client"

	"github.com/nats-io/nats.go"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

func buildDependencies(logger logger.Logger, consul *consulClient.ConsulClient) *do.Injector {
	injector := do.New()

	spotifySettingsValue, err := consul.Get("SpotifySettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain Spotify settings from Consul: %s", err.Error())
	}

	spotifySettings := spotifySettingsValue.(map[string]interface{})
	do.ProvideValue(injector, &accessToken.SpotifyClientConfig{
		ClientId:     spotifySettings["ClientId"].(string),
		ClientSecret: spotifySettings["ClientSecret"].(string),
	})

	natsSettingsValues, err := consul.Get("NatsSettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain Nats settings from Consul: '%s'", err.Error())
		return nil
	}
	natsSettings := natsSettingsValues.(map[string]interface{})

	natsConnection, err := nats.Connect(natsSettings["Url"].(string))
	if err != nil {
		logger.LogFatal(err, "Nats connection error: %s", err.Error())
	}

	do.ProvideValue(injector, natsConnection)

	do.Provide(injector, loggerServices.New)
	do.Provide(injector, common.ConstructHashCoder)
	do.Provide(injector, cache.ConstructMemoryCacheService)

	do.Provide(injector, labelServices.ConstructLabelService)
	do.Provide(injector, labelServices.ConstructManagerService)

	do.Provide(injector, spotifyServices.ConstructSpotifyService)

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
