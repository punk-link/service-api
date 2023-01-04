package api

import (
	"main/services/artists"
	"main/services/labels"
	"strconv"

	base "main/controllers"

	templates "github.com/punk-link/gin-generic-http-templates"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type ArtistController struct {
	artistService  *artists.ArtistService
	managerService *labels.ManagerService
}

func NewArtistController(injector *do.Injector) (*ArtistController, error) {
	artistService := do.MustInvoke[*artists.ArtistService](injector)
	managerService := do.MustInvoke[*labels.ManagerService](injector)

	return &ArtistController{
		artistService:  artistService,
		managerService: managerService,
	}, nil
}

func (t *ArtistController) Add(ctx *gin.Context) {
	spotifyId := ctx.Param("spotify-id")

	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		templates.NotFound(ctx, err.Error())
		return
	}

	result, err := t.artistService.Add(currentManager, spotifyId)
	templates.OkOrBadRequest(ctx, result, err)
}

func (t *ArtistController) Get(ctx *gin.Context) {
	labelId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		templates.BadRequest(ctx, err.Error())
		return
	}

	if err != nil {
		templates.NotFound(ctx, err.Error())
		return
	}

	result, err := t.artistService.Get(labelId)
	templates.OkOrBadRequest(ctx, result, err)
}

func (t *ArtistController) GetOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("artist-id"))
	if err != nil {
		templates.BadRequest(ctx, err.Error())
		return
	}

	if err != nil {
		templates.NotFound(ctx, err.Error())
		return
	}

	result, err := t.artistService.GetOne(id)
	templates.OkOrBadRequest(ctx, result, err)
}

func (t *ArtistController) Search(ctx *gin.Context) {
	query := ctx.Query("query")

	result := t.artistService.Search(query)
	templates.Ok(ctx, result)
}
