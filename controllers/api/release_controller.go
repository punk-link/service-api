package api

import (
	"main/services/artists"
	"strconv"

	"github.com/gin-gonic/gin"
	templates "github.com/punk-link/gin-generic-http-templates"
	"github.com/samber/do"
)

type ReleaseController struct {
	releaseService artists.ReleaseServer
}

func NewReleaseController(injector *do.Injector) (*ReleaseController, error) {
	releaseService := do.MustInvoke[artists.ReleaseServer](injector)

	return &ReleaseController{
		releaseService: releaseService,
	}, nil
}

func (t *ReleaseController) Get(ctx *gin.Context) {
	artistId, err := strconv.Atoi(ctx.Param("artist-id"))
	if err != nil {
		templates.BadRequest(ctx, err.Error())
		return
	}

	result, err := t.releaseService.Get(artistId)
	templates.OkOrBadRequest(ctx, result, err)
}

func (t *ReleaseController) GetOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		templates.BadRequest(ctx, err.Error())
		return
	}

	result, err := t.releaseService.GetOne(id)
	templates.OkOrBadRequest(ctx, result, err)
}
