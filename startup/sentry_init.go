package startup

import (
	"main/services/common"
	"main/utils"

	"github.com/getsentry/sentry-go"
)

func initSentry(logger *common.Logger) {
	err := sentry.Init(sentry.ClientOptions{
		AttachStacktrace: true,
		Dsn:              utils.GetEnvironmentVariable("SENTRY_DSN"),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		logger.LogError(err, "Sentry initialization failed: %v", err.Error())
	}
}
