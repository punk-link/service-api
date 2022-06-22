package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouters(app *gin.Engine) {
	v1 := app.Group("/v1")
	{
		v1.GET("/artist/search", controllers.SearchArtist)
		v1.GET("/artist/:spotify-id/releases", controllers.GetArtistReleases)
		v1.GET("/status", controllers.CheckStatus)
	}
}
