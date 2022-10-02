package mvc

import (
	base "main/controllers"
	"main/services/artists/mvc"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type MvcArtistController struct {
	service *mvc.MvcArtistService
}

func ConstructMvcArtistController(injector *do.Injector) (*MvcArtistController, error) {
	service := do.MustInvoke[*mvc.MvcArtistService](injector)

	return &MvcArtistController{
		service: service,
	}, nil
}

func (t *MvcArtistController) Get(ctx *gin.Context) {
	hash := ctx.Param("hash")

	result, err := t.service.Get(hash)
	base.OkOrNotFoundTemplate(ctx, "artist.tmpl", result, err)
}
