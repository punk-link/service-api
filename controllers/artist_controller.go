package controllers

import (
	"main/services/artists"

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

func (t *ArtistController) AddArtist(ctx *gin.Context) {
	spotifyId := ctx.Param("spotify-id")

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := t.artistService.AddArtist(currentManager, spotifyId)
	OkOrBadRequest(ctx, result, err)
}

func (t *ArtistController) SearchArtist(ctx *gin.Context) {
	query := ctx.Query("query")

	result := t.artistService.SearchArtist(query)
	Ok(ctx, result)
}
