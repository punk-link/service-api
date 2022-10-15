package api

import (
	base "main/controllers"
	"main/data"

	"github.com/gin-gonic/gin"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type StatusController struct {
	logger *logger.Logger
}

func ConstructStatusController(injector *do.Injector) (*StatusController, error) {
	logger := do.MustInvoke[*logger.Logger](injector)

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
