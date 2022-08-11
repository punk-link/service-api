package controllers

import (
	"main/services/artists"
	"main/services/spotify"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArtistController struct {
	artistService  *artists.ArtistService
	spotifyService *spotify.SpotifyService
}

func BuildArtistController(artistService *artists.ArtistService, spotifyService *spotify.SpotifyService) *ArtistController {
	return &ArtistController{
		artistService:  artistService,
		spotifyService: spotifyService,
	}
}

func (t *ArtistController) GetRelease(ctx *gin.Context) {
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

	result := t.artistService.GetRelease(currentManager, id)
	Ok(ctx, result)
}

func (t *ArtistController) GetReleases(ctx *gin.Context) {
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

	result := t.artistService.GetReleases(currentManager, artistId)
	Ok(ctx, result)
}

func (t *ArtistController) SearchArtist(ctx *gin.Context) {
	query := ctx.Query("query")

	result := t.artistService.SearchArtist(query)
	Ok(ctx, result)
}
