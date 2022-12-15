package startup

import (
	controllers "main/controllers"
	apiControllers "main/controllers/api"
	staticControllers "main/controllers/static"
	dataConfig "main/data"
	artistModels "main/models/artists"
	"main/models/platforms/spotify/tokens"
	artistServices "main/services/artists"
	artistStaticServices "main/services/artists/static"
	"main/services/common"
	labelServices "main/services/labels"
	platformServices "main/services/platforms"
	spotifyServices "main/services/platforms/spotify"

	"github.com/nats-io/nats.go"
	cacheManager "github.com/punk-link/cache-manager"
	consulClient "github.com/punk-link/consul-client"
	loggerService "github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func buildDependencies(logger loggerService.Logger, consul consulClient.ConsulClient /*, appSecrets map[string]any*/) *do.Injector {
	injector := do.New()

	spotifySettingsValue, err := consul.Get("SpotifySettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain Spotify settings from Consul: %s", err.Error())
	}

	spotifySettings := spotifySettingsValue.(map[string]any)
	do.ProvideValue(injector, &tokens.SpotifyClientConfig{
		ClientId:     spotifySettings["ClientId"].(string),
		ClientSecret: spotifySettings["ClientSecret"].(string), //appSecrets["client-secret"].(string),
	})

	natsSettingsValues, err := consul.Get("NatsSettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain Nats settings from Consul: '%s'", err.Error())
		return nil
	}
	natsSettings := natsSettingsValues.(map[string]any)

	natsConnection, err := nats.Connect(natsSettings["Endpoint"].(string))
	if err != nil {
		logger.LogFatal(err, "Nats connection error: %s", err.Error())
	}

	do.ProvideValue(injector, natsConnection)

	do.Provide(injector, func(i *do.Injector) (loggerService.Logger, error) {
		return loggerService.New(), nil
	})
	do.Provide(injector, func(i *do.Injector) (*gorm.DB, error) {
		return dataConfig.New(logger, consul /*, appSecrets*/), nil
	})

	do.Provide(injector, registerCacheManager[artistModels.Artist]())
	do.Provide(injector, registerCacheManager[artistModels.Release]())
	do.Provide(injector, registerCacheManager[[]artistModels.Release]())
	do.Provide(injector, registerCacheManager[map[string]any]())

	do.Provide(injector, common.NewHashCoder)
	do.Provide(injector, labelServices.NewLabelService)
	do.Provide(injector, labelServices.NewManagerService)

	do.Provide(injector, spotifyServices.NewSpotifyService)

	do.Provide(injector, artistServices.NewReleaseService)
	do.Provide(injector, artistStaticServices.NewStaticReleaseService)
	do.Provide(injector, artistServices.NewArtistService)
	do.Provide(injector, artistStaticServices.NewStaticArtistService)

	do.Provide(injector, platformServices.NewStreamingPlatformService)

	do.Provide(injector, controllers.NewMetricsController)
	do.Provide(injector, controllers.NewStatusController)

	do.Provide(injector, apiControllers.NewArtistController)
	do.Provide(injector, apiControllers.NewHashController)
	do.Provide(injector, apiControllers.NewLabelController)
	do.Provide(injector, apiControllers.NewManagerController)
	do.Provide(injector, apiControllers.NewStreamingPlatformController)
	do.Provide(injector, apiControllers.NewReleaseController)

	do.Provide(injector, staticControllers.NewStaticArtistController)
	do.Provide(injector, staticControllers.NewStaticReleaseController)

	return injector
}

func registerCacheManager[T any]() func(i *do.Injector) (cacheManager.CacheManager[T], error) {
	return func(i *do.Injector) (cacheManager.CacheManager[T], error) {
		return cacheManager.New[T](), nil
	}
}
