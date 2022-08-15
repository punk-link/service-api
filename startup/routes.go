package startup

import (
	"main/controllers"
	"main/services/common"
	"main/services/helpers"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func setupRouts(app *gin.Engine, diContainer *dig.Container, logger *common.Logger) {
	registerRoutes(logger, diContainer, controllers.StatusController{}, func(controller *controllers.StatusController) {
		app.GET("/health", controller.CheckHealth)
	})

	v1 := app.Group("/v1")

	registerRoutes(logger, diContainer, controllers.ManagerController{}, func(controller *controllers.ManagerController) {
		v1.POST("/managers", controller.AddManager)
		v1.POST("/managers/master", controller.AddMasterManager)
		v1.GET("/managers", controller.GetManagers)
		v1.GET("/managers/:id", controller.GetManager)
		v1.POST("/managers/:id", controller.ModifyManager)
	})

	registerRoutes(logger, diContainer, controllers.LabelController{}, func(controller *controllers.LabelController) {
		v1.GET("/labels/:id", controller.GetLabel)
		v1.POST("/labels/:id", controller.ModifyLabel)
	})

	registerRoutes(logger, diContainer, controllers.ArtistController{}, func(controller *controllers.ArtistController) {
		v1.POST("/artists/:spotify-id", controller.AddArtist)
		v1.GET("/artists/search", controller.SearchArtist)
	})

	registerRoutes(logger, diContainer, controllers.ReleaseController{}, func(controller *controllers.ReleaseController) {
		v1.GET("/artists/:artist-id/releases", controller.GetReleases)
		v1.GET("/artists/releases/:id/", controller.GetRelease)
	})
}

func registerRoutes[T any](logger *common.Logger, diContainer *dig.Container, target T, function func(*T)) {
	err := diContainer.Invoke(function)
	if err != nil {
		structName := helpers.GetStructNameAsString(target)
		logger.LogFatal(err, "Can't resolve a dependency '%s': %v", structName, err.Error())
		panic(err.Error())
	}
}
