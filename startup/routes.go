package startup

import (
	apiControllers "main/controllers/api"
	staticControllers "main/controllers/static"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func setupRouts(app *gin.Engine, diContainer *do.Injector) {
	registerRoutes(diContainer, func(controller *apiControllers.StatusController) {
		app.GET("/health", controller.CheckHealth)
	})

	registerRoutes(diContainer, func(controller *staticControllers.StaticArtistController) {
		app.GET("/artists/:hash", controller.Get)
	})

	registerRoutes(diContainer, func(controller *staticControllers.StaticReleaseController) {
		app.GET("/releases/:hash", controller.Get)
	})

	v1 := app.Group("/v1")

	registerRoutes(diContainer, func(controller *apiControllers.HashController) {
		v1.GET("/hashes/:target/decode", controller.Decode)
		v1.GET("/hashes/:target/encode", controller.Encode)
	})

	registerRoutes(diContainer, func(controller *apiControllers.ManagerController) {
		v1.POST("/managers", controller.Add)
		v1.POST("/managers/master", controller.AddMaster)
		v1.GET("/managers", controller.Get)
		v1.GET("/managers/:id", controller.GetOne)
		v1.POST("/managers/:id", controller.Modify)
	})

	registerRoutes(diContainer, func(controller *apiControllers.LabelController) {
		subroutes := v1.Group("/labels")
		{
			subroutes.GET("/:id", controller.Get)
			subroutes.POST("/:id", controller.Modify)
			registerRoutes(diContainer, func(controller *apiControllers.ArtistController) {
				subroutes.GET("/:id/artists/", controller.Get)
			})
		}
	})

	registerRoutes(diContainer, func(controller *apiControllers.ArtistController) {
		v1.POST("/artists/:spotify-id", controller.Add)
		v1.GET("/artists/search", controller.Search)
		subroutes := v1.Group("/artists")
		{
			subroutes.GET("/:artist-id", controller.GetOne)
			registerRoutes(diContainer, func(controller *apiControllers.ReleaseController) {
				subroutes.GET("/:artist-id/releases/", controller.Get)
				subroutes.GET("/:artist-id/releases/:id/", controller.GetOne)
			})
		}
	})

	registerRoutes(diContainer, func(controller *apiControllers.PlatformSynchronisationController) {
		v1.POST("/platforms/sync", controller.Sync)
	})
}

func registerRoutes[T any](diContainer *do.Injector, function func(*T)) {
	controller := do.MustInvoke[*T](diContainer)
	function(controller)
}
