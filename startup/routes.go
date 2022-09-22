package startup

import (
	apiControllers "main/controllers/api"
	mvcControllers "main/controllers/mvc"
	"main/helpers"
	"main/services/common"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func setupRouts(app *gin.Engine, diContainer *dig.Container, logger *common.Logger) {
	registerRoutes(logger, diContainer, func(controller *apiControllers.StatusController) {
		app.GET("/health", controller.CheckHealth)
	})

	registerRoutes(logger, diContainer, func(controller *mvcControllers.MvcReleaseController) {
		app.GET("/releases/:hash", controller.Get)
	})

	v1 := app.Group("/v1")

	registerRoutes(logger, diContainer, func(controller *apiControllers.HashController) {
		v1.GET("/hashes/:target/decode", controller.Decode)
		v1.GET("/hashes/:target/encode", controller.Encode)
	})

	registerRoutes(logger, diContainer, func(controller *apiControllers.ManagerController) {
		v1.POST("/managers", controller.Add)
		v1.POST("/managers/master", controller.AddMaster)
		v1.GET("/managers", controller.Get)
		v1.GET("/managers/:id", controller.GetOne)
		v1.POST("/managers/:id", controller.Modify)
	})

	registerRoutes(logger, diContainer, func(controller *apiControllers.LabelController) {
		subroutes := v1.Group("/labels")
		{
			subroutes.GET("/:id", controller.Get)
			subroutes.POST("/:id", controller.Modify)
			registerRoutes(logger, diContainer, func(controller *apiControllers.ArtistController) {
				subroutes.GET("/:id/artists/", controller.Get)
			})
		}
	})

	registerRoutes(logger, diContainer, func(controller *apiControllers.ArtistController) {
		v1.POST("/artists/:spotify-id", controller.Add)
		v1.GET("/artists/search", controller.Search)
		subroutes := v1.Group("/artists")
		{
			subroutes.GET("/:artist-id", controller.GetOne)
			registerRoutes(logger, diContainer, func(controller *apiControllers.ReleaseController) {
				subroutes.GET("/:artist-id/releases/", controller.Get)
				subroutes.GET("/:artist-id/releases/:id/", controller.GetOne)
			})
		}
	})
}

func registerRoutes[T any](logger *common.Logger, diContainer *dig.Container, function func(*T)) {
	err := diContainer.Invoke(function)
	if err != nil {
		target := new(T)
		structName := helpers.GetStructNameAsString(target)
		logger.LogFatal(err, "Can't resolve a dependency '%s': %v", structName, err.Error())
		panic(err.Error())
	}
}
