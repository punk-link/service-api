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

	initSentry(app, logger, consul)
	setupRouts(app, diContainer, logger)

	data.ConfigureDatabase(logger, consul)

	return app
}
