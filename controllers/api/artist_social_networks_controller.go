package api

import (
	base "main/controllers"
	artistModels "main/models/artists"
	artistServices "main/services/artists"
	labelServices "main/services/labels"
	"strconv"

	templates "github.com/punk-link/gin-generic-http-templates"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type ArtistSocialNetworksController struct {
	managerService       labelServices.ManagerServer
	socialNetworkService artistServices.SocialNetworkServer
}

func NewArtistSocialNetworksController(injector *do.Injector) (*ArtistSocialNetworksController, error) {
	managerService := do.MustInvoke[labelServices.ManagerServer](injector)
	socialNetworkService := do.MustInvoke[artistServices.SocialNetworkServer](injector)

	return &ArtistSocialNetworksController{
		managerService:       managerService,
		socialNetworkService: socialNetworkService,
	}, nil
}

func (t *ArtistSocialNetworksController) AddOrModify(ctx *gin.Context) {
	artistId, err := strconv.Atoi(ctx.Param("artist-id"))
	if err != nil {
		templates.BadRequest(ctx, err.Error())
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		templates.NotFound(ctx, err.Error())
		return
	}

	var networks []artistModels.SocialNetwork
	if err := ctx.ShouldBindJSON(&networks); err != nil {
		templates.UnprocessableEntity(ctx, err)
		return
	}

	results, err := t.socialNetworkService.ArrOrModify(currentManager, artistId, networks)

	templates.OkOrBadRequest(ctx, results, err)
}

func (t *ArtistSocialNetworksController) Get(ctx *gin.Context) {
	artistId, err := strconv.Atoi(ctx.Param("artist-id"))
	if err != nil {
		templates.BadRequest(ctx, err.Error())
		return
	}

	result := t.socialNetworkService.Get(artistId)
	templates.OkOrBadRequest(ctx, result, nil)
}
