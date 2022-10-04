package startup

import (
	"main/data"
	"main/infrastructure/consul"
	"main/services/common"

	"github.com/gin-gonic/gin"
)

func Configure(logger *common.Logger, consul *consul.ConsulClient, ginMode string) *gin.Engine {
	diContainer := buildDependencies()

	gin.SetMode(ginMode)
	app := gin.Default()

	app.LoadHTMLGlob("./var/www/templates/**/*.go.tmpl")
	app.Static("/assets", "./var/www/assets")

	initSentry(app, logger, consul)
	setupRouts(app, diContainer)

	data.ConfigureDatabase(logger, consul)

	return app
}
