package api

import (
	base "main/controllers"
	"main/services/artists"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArtistController struct {
	artistService *artists.ArtistService
}

func ConstructArtistController(artistService *artists.ArtistService) *ArtistController {
	return &ArtistController{
		artistService: artistService,
	}
}

func (t *ArtistController) Add(ctx *gin.Context) {
	spotifyId := ctx.Param("spotify-id")

	currentManager, err := base.GetCurrentManagerContext(ctx)
	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := t.artistService.Add(currentManager, spotifyId)
	base.OkOrBadRequest(ctx, result, err)
}

func (t *ArtistController) Get(ctx *gin.Context) {
	labelId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := t.artistService.Get(labelId)
	base.OkOrBadRequest(ctx, result, err)
}

func (t *ArtistController) GetOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("artist-id"))
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := t.artistService.GetOne(id)
	base.OkOrBadRequest(ctx, result, err)
}

func (t *ArtistController) Search(ctx *gin.Context) {
	query := ctx.Query("query")

	result := t.artistService.Search(query)
	base.Ok(ctx, result)
}
