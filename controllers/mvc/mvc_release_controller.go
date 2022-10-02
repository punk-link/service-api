package mvc

import (
	base "main/controllers"
	"main/services/artists/mvc"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type MvcReleaseController struct {
	service *mvc.MvcReleaseService
}

func ConstructMvcReleaseController(injector *do.Injector) (*MvcReleaseController, error) {
	service := do.MustInvoke[*mvc.MvcReleaseService](injector)

	return &MvcReleaseController{
		service: service,
	}, nil
}

func (t *MvcReleaseController) Get(ctx *gin.Context) {
	hash := ctx.Param("hash")

	result, err := t.service.Get(hash)
	base.OkOrNotFoundTemplate(ctx, "release.tmpl", result, err)
}
