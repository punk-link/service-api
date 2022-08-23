package startup

import (
	"main/controllers"
	"main/helpers"
	"main/services/common"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func setupRouts(app *gin.Engine, diContainer *dig.Container, logger *common.Logger) {
	registerRoutes(logger, diContainer, controllers.StatusController{}, func(controller *controllers.StatusController) {
		app.GET("/health", controller.CheckHealth)
	})

	v1 := app.Group("/v1")

	registerRoutes(logger, diContainer, controllers.ManagerController{}, func(controller *controllers.ManagerController) {
		v1.POST("/managers", controller.Add)
		v1.POST("/managers/master", controller.AddMaster)
		v1.GET("/managers", controller.Get)
		v1.GET("/managers/:id", controller.GetOne)
		v1.POST("/managers/:id", controller.Modify)
	})

	registerRoutes(logger, diContainer, controllers.LabelController{}, func(controller *controllers.LabelController) {
		v1.GET("/labels/:id", controller.Get)
		v1.POST("/labels/:id", controller.Modify)
	})

	registerRoutes(logger, diContainer, controllers.ArtistController{}, func(controller *controllers.ArtistController) {
		v1.POST("/artists/:spotify-id", controller.Add)
		v1.GET("/artists/:id", controller.Get)
		v1.GET("/artists/search", controller.Search)
	})

	registerRoutes(logger, diContainer, controllers.ReleaseController{}, func(controller *controllers.ReleaseController) {
		v1.GET("/artists/:artist-id/releases", controller.Get)
		v1.GET("/artists/releases/:id/", controller.GetOne)
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
