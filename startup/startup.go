package startup

import (
	"main/data"

	"github.com/gin-gonic/gin"
	consulClient "github.com/punk-link/consul-client"
	"github.com/punk-link/logger"
)

func Configure(logger logger.Logger, consul *consulClient.ConsulClient, ginMode string) *gin.Engine {
	diContainer := buildDependencies(logger, consul)

	gin.SetMode(ginMode)
	app := gin.Default()

	app.LoadHTMLGlob("./var/www/templates/**/*.go.tmpl")
	app.Static("/assets", "./var/www/assets")

	initSentry(app, logger, consul)
	setupRouts(app, diContainer)

	data.ConfigureDatabase(logger, consul)

	return app
}
