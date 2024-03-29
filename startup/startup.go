package startup

import (
	startupModels "main/models/startup"

	"github.com/gin-gonic/gin"
	consulClient "github.com/punk-link/consul-client"
	"github.com/punk-link/logger"
	"github.com/samber/do"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Configure(logger logger.Logger, consul consulClient.ConsulClient, injector *do.Injector, options *startupModels.StartupOptions) *gin.Engine {
	gin.SetMode(options.GinMode)
	app := gin.Default()

	app.Use(metricsMiddleware(options.ServiceName))
	app.Use(otelgin.Middleware(options.ServiceName))

	initSentry(app, logger, consul, options.EnvironmentName)
	configureOpenTelemetry(logger, consul, options)
	setupRouts(app, injector)

	return app
}
