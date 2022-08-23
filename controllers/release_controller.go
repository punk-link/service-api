package controllers

import (
	"main/services/artists"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReleaseController struct {
	releaseService *artists.ReleaseService
}

func BuildReleaseController(releaseService *artists.ReleaseService) *ReleaseController {
	return &ReleaseController{
		releaseService: releaseService,
	}
}

func (t *ReleaseController) GetRelease(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result := t.releaseService.GetRelease(currentManager, id)
	Ok(ctx, result)
}

func (t *ReleaseController) GetReleases(ctx *gin.Context) {
	artistId, err := strconv.Atoi(ctx.Param("artist-id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := t.releaseService.GetReleases(currentManager, artistId)
	OkOrBadRequest(ctx, result, err)
}
