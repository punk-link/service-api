package routers

import (
	"main/controllers"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

func SetupRouters(app *gin.Engine, diContainer *dig.Container) {
	app.GET("/health", controllers.CheckHealth)

	err := diContainer.Invoke(func(managerController *controllers.ManagerController) {
		v1 := app.Group("/v1")
		{
			v1.POST("/managers", managerController.AddManager)
			v1.POST("/managers/master", managerController.AddMasterManager)
			v1.GET("/managers", managerController.GetManagers)
			v1.GET("/managers/:id", managerController.GetManager)
			v1.POST("/managers/:id", managerController.ModifyManager)
		}
	})
	if err != nil {
		log.Error().Err(err).Msgf("Can't resolve a dependency '%s': %v", reflect.TypeOf(controllers.ManagerController{}).Name(), err.Error())
		panic(err.Error())
	}

	err = diContainer.Invoke(func(labelController *controllers.LabelController) {
		v1 := app.Group("/v1")
		{
			v1.GET("/labels/:id", labelController.GetLabel)
			v1.POST("/labels/:id", labelController.ModifyLabel)
		}
	})
	if err != nil {
		log.Error().Err(err).Msgf("Can't resolve a dependency '%s': %v", reflect.TypeOf(controllers.LabelController{}).Name(), err.Error())
		panic(err.Error())
	}

	v1 := app.Group("/v1")
	{
		v1.GET("/artists/search", controllers.SearchArtist)
		v1.GET("/artists/:spotify-id/releases", controllers.GetReleases)
		v1.GET("/artists/releases/:spotify-id/", controllers.GetRelease)
	}
}
