package mvc

import (
	base "main/controllers"
	"main/services/artists"

	"github.com/gin-gonic/gin"
)

type MvcReleaseController struct {
	service *artists.MvcReleaseService
}

func ConstructMvcReleaseController(service *artists.MvcReleaseService) *MvcReleaseController {
	return &MvcReleaseController{
		service: service,
	}
}

func (t *MvcReleaseController) Get(ctx *gin.Context) {
	hash := ctx.Param("hash")

	result, err := t.service.Get(hash)
	base.OkOrNotFoundTemplate(ctx, "release.tmpl", result, err)
}
