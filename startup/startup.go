package startup

import (
	"main/controllers"
	"main/routers"
	labelServices "main/services/labels"
	"main/utils"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

func Configure() *gin.Engine {
	err := sentry.Init(sentry.ClientOptions{
		AttachStacktrace: true,
		Dsn:              utils.GetEnvironmentVariable("SENTRY_DSN"),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Error().Err(err).Msgf("Sentry initialization failed: %v", err.Error())
	}

	diContainer := dig.New()
	diContainer.Provide(labelServices.NewLabelService)
	diContainer.Provide(labelServices.NewManagerService)
	diContainer.Provide(controllers.NewManagerController)
	diContainer.Provide(controllers.NewLabelController)

	app := gin.Default()
	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))
	routers.SetupRouters(app, diContainer)

	return app
}
