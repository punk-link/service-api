package api

import (
	base "main/controllers"

	"github.com/gin-gonic/gin"
	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type StatusController struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewStatusController(injector *do.Injector) (*StatusController, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	return &StatusController{
		db:     db,
		logger: logger,
	}, nil
}

func (t *StatusController) CheckHealth(ctx *gin.Context) {
	sqlDb, err := t.db.DB()
	if err != nil {
		t.logger.LogError(err, "Postgres initialization failed: %v", err.Error())
		base.InternalServerError(ctx, err.Error())
		return
	}

	err = sqlDb.Ping()
	if err != nil {
		t.logger.LogError(err, "Can't reach any database instances: %v", err.Error())
		base.InternalServerError(ctx, err.Error())
		return
	}

	base.Ok(ctx, "OK")
}
