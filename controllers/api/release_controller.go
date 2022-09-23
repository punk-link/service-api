package api

import (
	base "main/controllers"
	"main/services/artists"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReleaseController struct {
	releaseService *artists.ReleaseService
}

func ConstructReleaseController(releaseService *artists.ReleaseService) *ReleaseController {
	return &ReleaseController{
		releaseService: releaseService,
	}
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
