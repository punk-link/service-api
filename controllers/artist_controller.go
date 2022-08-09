package controllers

import (
	"main/services/spotify"

	"github.com/gin-gonic/gin"
)

type ArtistController struct {
	spotifyService *spotify.SpotifyService
}

func BuildArtistController(spotifyService *spotify.SpotifyService) *ArtistController {
	return &ArtistController{
		spotifyService: spotifyService,
	}
}

func (controller *ArtistController) GetRelease(ctx *gin.Context) {
	spotifyId := ctx.Param("spotify-id")

	result := controller.spotifyService.GetArtistRelease(spotifyId)
	Ok(ctx, result)
}

func (controller *ArtistController) GetReleases(ctx *gin.Context) {
	spotifyId := ctx.Param("spotify-id")

	result := controller.spotifyService.GetArtistReleases(spotifyId)
	Ok(ctx, result)
}

func (controller *ArtistController) SearchArtist(ctx *gin.Context) {
	query := ctx.Query("query")

	result := controller.spotifyService.SearchArtist(query)
	Ok(ctx, result)
}
