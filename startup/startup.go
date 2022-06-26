package startup

import (
	"main/routers"
	"main/utils"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

	app := gin.Default()
	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))
	routers.SetupRouters(app)

	return app
}
