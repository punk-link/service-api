package startup

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	consulClient "github.com/punk-link/consul-client"
	"github.com/punk-link/logger"
)

func initSentry(app *gin.Engine, logger logger.Logger, consul *consulClient.ConsulClient) {
	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	dsn, _ := consul.Get("SentryDsn")
	err := sentry.Init(sentry.ClientOptions{
		AttachStacktrace: true,
		Dsn:              dsn.(string),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		logger.LogError(err, "Sentry initialization failed: %v", err.Error())
	}
}
