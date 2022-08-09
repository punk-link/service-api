package startup

import (
	"main/controllers"
	"main/services/common"
	labelServices "main/services/labels"
	spotifyServices "main/services/spotify"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

func Configure() *gin.Engine {
	diContainer := dig.New()

	diContainer.Provide(common.BuildLogger)
	diContainer.Provide(labelServices.BuildLabelService)
	diContainer.Provide(labelServices.BuildManagerService)
	diContainer.Provide(spotifyServices.BuildSpotifyService)

	diContainer.Provide(controllers.BuildArtistController)
	diContainer.Provide(controllers.BuildLabelController)
	diContainer.Provide(controllers.BuildManagerController)
	diContainer.Provide(controllers.BuildStatusController)

	app := gin.Default()
	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	err := diContainer.Invoke(func(logger *common.Logger) {
		initSentry(logger)
		setupRouts(app, diContainer, logger)
	})
	if err != nil {
		log.Error().Err(err).Msgf("Can't resolve a logger: %v", err.Error())
	}

	return app
}
