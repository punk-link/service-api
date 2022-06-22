package controllers

import (
	"main/services/spotify"

	"github.com/gin-gonic/gin"
)

func SearchArtist(ctx *gin.Context) {
	query := ctx.Query("query")

	result, err := spotify.SearchArtist(query)
	OkOrBadRequest(ctx, result, err)
}
