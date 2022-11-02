package static

import (
	base "main/controllers"
	"main/services/artists/static"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type StaticArtistController struct {
	service *static.StaticArtistService
}

func NewStaticArtistController(injector *do.Injector) (*StaticArtistController, error) {
	service := do.MustInvoke[*static.StaticArtistService](injector)

	return &StaticArtistController{
		service: service,
	}, nil
}

func (t *StaticArtistController) Get(ctx *gin.Context) {
	hash := ctx.Param("hash")

	result, err := t.service.Get(hash)
	base.OkOrNotFoundTemplate(ctx, "artist.go.tmpl", result, err)
}
