package api

import (
	base "main/controllers"
	"main/data"
	"main/services/common"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type StatusController struct {
	logger *common.Logger
}

func ConstructStatusController(injector *do.Injector) (*StatusController, error) {
	logger := do.MustInvoke[*common.Logger](injector)

	return &StatusController{
		logger: logger,
	}, nil
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
