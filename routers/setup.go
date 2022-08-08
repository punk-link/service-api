package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouters(app *gin.Engine) {
	app.GET("/health", controllers.CheckHealth)

	v1 := app.Group("/v1")
	{
		v1.GET("/artists/search", controllers.SearchArtist)
		v1.GET("/artists/:spotify-id/releases", controllers.GetReleases)
		v1.GET("/artists/releases/:spotify-id/", controllers.GetRelease)

		v1.POST("/managers", controllers.AddManager)
		v1.POST("/managers/master", controllers.AddMasterManager)
		v1.GET("/managers", controllers.GetManagers)
		v1.GET("/managers/:id", controllers.GetManager)
		v1.POST("/managers/:id", controllers.ModifyManager)

		v1.GET("/labels/:id", controllers.GetLabel)
		v1.POST("/labels/:id", controllers.ModifyLabel)
	}
}
