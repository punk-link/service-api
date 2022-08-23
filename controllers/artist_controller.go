package controllers

import (
	"main/services/artists"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArtistController struct {
	artistService *artists.ArtistService
}

func BuildArtistController(artistService *artists.ArtistService) *ArtistController {
	return &ArtistController{
		artistService: artistService,
	}
}

func (t *ArtistController) Add(ctx *gin.Context) {
	spotifyId := ctx.Param("spotify-id")

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := t.artistService.Add(currentManager, spotifyId)
	OkOrBadRequest(ctx, result, err)
}

func (t *ArtistController) Get(ctx *gin.Context) {
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

	result, err := t.artistService.Get(currentManager, id)
	OkOrBadRequest(ctx, result, err)
}

func (t *ArtistController) Search(ctx *gin.Context) {
	query := ctx.Query("query")

	result := t.artistService.Search(query)
	Ok(ctx, result)
}
