package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouters(app *gin.Engine) {
	v1 := app.Group("/v1")
	{
		v1.GET("/artists/search", controllers.SearchArtist)
		v1.GET("/artists/:spotify-id/releases", controllers.GetReleases)
		v1.GET("/artists/releases/:spotify-id/", controllers.GetRelease)

		v1.POST("/managers", controllers.AddManager)
		v1.POST("/managers/master", controllers.AddMasterManager)
		v1.GET("/managers/:id", controllers.GetManager)
		v1.POST("/managers/:id", controllers.ModifyManager)

		v1.GET("/organizations/:id", controllers.GetOrganization)
		v1.POST("/organizations/:id", controllers.ModifyOrganization)

		v1.GET("/status", controllers.CheckStatus)
	}
}
