package startup

import (
	"main/infrastructure/consul"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/punk-link/logger"
)

func initSentry(app *gin.Engine, logger *logger.Logger, consul *consul.ConsulClient) {
	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	dsn := consul.Get("SentryDsn").(string)
	err := sentry.Init(sentry.ClientOptions{
		AttachStacktrace: true,
		Dsn:              dsn,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		logger.LogError(err, "Sentry initialization failed: %v", err.Error())
	}
}
