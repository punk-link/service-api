package mvc

import (
	base "main/controllers"
	"main/services/artists"

	"github.com/gin-gonic/gin"
)

type MvcArtistController struct {
	service *artists.MvcArtistService
}

func ConstructMvcArtistController(service *artists.MvcArtistService) *MvcArtistController {
	return &MvcArtistController{
		service: service,
	}
}

func (t *MvcArtistController) Get(ctx *gin.Context) {
	hash := ctx.Param("hash")

	result, err := t.service.Get(hash)
	base.OkOrNotFoundTemplate(ctx, "artist.tmpl", result, err)
}
