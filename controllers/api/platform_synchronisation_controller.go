package api

import (
	base "main/controllers"
	"main/services/platforms"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type PlatformSynchronisationController struct {
	platformSynchronisationService *platforms.StrimingPlatformService
}

func ConstructPlatformSynchronisationController(injector *do.Injector) (*PlatformSynchronisationController, error) {
	platformSynchronisationService := do.MustInvoke[*platforms.StrimingPlatformService](injector)

	return &PlatformSynchronisationController{
		platformSynchronisationService: platformSynchronisationService,
	}, nil
}

func (t *PlatformSynchronisationController) Sync(ctx *gin.Context) {
	t.platformSynchronisationService.SyncUrls()

	base.NoContent(ctx)
}
