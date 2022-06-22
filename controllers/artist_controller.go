package controllers

import (
	"main/services/spotify"

	"github.com/gin-gonic/gin"
)

func GetArtistReleases(ctx *gin.Context) {
	spotifyId := ctx.Param("spotify-id")

	result, err := spotify.GetArtistReleases(spotifyId)
	OkOrBadRequest(ctx, result, err)
}

func SearchArtist(ctx *gin.Context) {
	query := ctx.Query("query")

	result, err := spotify.SearchArtist(query)
	OkOrBadRequest(ctx, result, err)
}
