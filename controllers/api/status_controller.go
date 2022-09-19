package api

import (
	base "main/controllers"
	"main/data"
	"main/services/common"

	"github.com/gin-gonic/gin"
)

type StatusController struct {
	logger *common.Logger
}

func ConstructStatusController(logger *common.Logger) *StatusController {
	return &StatusController{
		logger: logger,
	}
}

func (controller *StatusController) CheckHealth(ctx *gin.Context) {
	sqlDb, err := data.DB.DB()
	if err != nil {
		controller.logger.LogError(err, "Postgres initialization failed: %v", err.Error())
		base.InternalServerError(ctx, err.Error())
		return
	}

	err = sqlDb.Ping()
	if err != nil {
		controller.logger.LogError(err, "Can't reach any database instances: %v", err.Error())
		base.InternalServerError(ctx, err.Error())
		return
	}

	base.Ok(ctx, "OK")
}
