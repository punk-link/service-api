package api

import (
	base "main/controllers"
	"main/services/artists"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type ReleaseController struct {
	releaseService *artists.ReleaseService
}

func NewReleaseController(injector *do.Injector) (*ReleaseController, error) {
	releaseService := do.MustInvoke[*artists.ReleaseService](injector)

	return &ReleaseController{
		releaseService: releaseService,
	}, nil
}

func (t *ReleaseController) GetOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	result, err := t.releaseService.GetOne(id)
	base.OkOrBadRequest(ctx, result, err)
}

func (t *ReleaseController) Get(ctx *gin.Context) {
	artistId, err := strconv.Atoi(ctx.Param("artist-id"))
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	result, err := t.releaseService.GetByArtistId(artistId)
	base.OkOrBadRequest(ctx, result, err)
}
