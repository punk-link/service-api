package controllers

import (
	"main/services/spotify"

	"github.com/gin-gonic/gin"
)

func GetRelease(ctx *gin.Context) {
	spotifyId := ctx.Param("spotify-id")

	result := spotify.GetArtistRelease(spotifyId)
	Ok(ctx, result)
}

func GetReleases(ctx *gin.Context) {
	spotifyId := ctx.Param("spotify-id")

	result := spotify.GetArtistReleases(spotifyId)
	Ok(ctx, result)
}

func SearchArtist(ctx *gin.Context) {
	query := ctx.Query("query")

	result := spotify.SearchArtist(query)
	Ok(ctx, result)
}
