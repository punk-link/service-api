package static

import (
	base "main/controllers"
	"main/services/artists/static"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type StaticReleaseController struct {
	service *static.StaticReleaseService
}

func ConstructStaticReleaseController(injector *do.Injector) (*StaticReleaseController, error) {
	service := do.MustInvoke[*static.StaticReleaseService](injector)

	return &StaticReleaseController{
		service: service,
	}, nil
}

func (t *StaticReleaseController) Get(ctx *gin.Context) {
	hash := ctx.Param("hash")

	result, err := t.service.Get(hash)
	base.OkOrNotFoundTemplate(ctx, "release.go.tmpl", result, err)
}
