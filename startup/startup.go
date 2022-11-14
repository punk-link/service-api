package startup

import (
	"main/data"
	startupModels "main/models/startup"

	"github.com/gin-gonic/gin"
	consulClient "github.com/punk-link/consul-client"
	"github.com/punk-link/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Configure(logger logger.Logger, consul *consulClient.ConsulClient, options *startupModels.StartupOptions) *gin.Engine {
	diContainer := buildDependencies(logger, consul)

	gin.SetMode(options.GinMode)
	app := gin.Default()

	app.Use(metricsMiddleware(options.ServiceName))
	app.Use(otelgin.Middleware(options.ServiceName))

	app.LoadHTMLGlob("./var/www/templates/**/*.go.tmpl")
	app.Static("/assets", "./var/www/assets")

	initSentry(app, logger, consul, options.EnvironmentName)
	configureOpenTelemetry(logger, options)
	setupRouts(app, diContainer)

	data.New(logger, consul)

	return app
}
