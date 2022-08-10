package startup

import (
	"main/infrastructure"
	"main/services/common"

	"github.com/getsentry/sentry-go"
)

func initSentry(logger *common.Logger) {
	err := sentry.Init(sentry.ClientOptions{
		AttachStacktrace: true,
		Dsn:              infrastructure.GetEnvironmentVariable("SENTRY_DSN"),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		logger.LogError(err, "Sentry initialization failed: %v", err.Error())
	}
}
