package startup

import (
	"main/controllers"
	apiControllers "main/controllers/api"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func setupRouts(app *gin.Engine, injector *do.Injector) {
	registerRoutes(injector, func(controller *controllers.MetricsController) {
		app.GET("/metrics", controller.GetMetrics)
	})

	registerRoutes(injector, func(controller *controllers.StatusController) {
		app.GET("/health", controller.CheckHealth)
		app.GET("/error", controller.ThrowError)
	})

	v1 := app.Group("/v1")

	registerRoutes(injector, func(controller *apiControllers.ManagerController) {
		v1.POST("/managers", controller.Add)
		v1.POST("/managers/master", controller.AddMaster)
		v1.GET("/managers", controller.Get)
		v1.GET("/managers/:id", controller.GetOne)
		v1.POST("/managers/:id", controller.Modify)
	})

	registerRoutes(injector, func(controller *apiControllers.LabelController) {
		subroutes := v1.Group("/labels")
		{
			subroutes.GET("/:id", controller.Get)
			subroutes.POST("/:id", controller.Modify)
			registerRoutes(injector, func(controller *apiControllers.ArtistController) {
				subroutes.GET("/:id/artists/", controller.Get)
			})
		}
	})

	registerRoutes(injector, func(controller *apiControllers.ArtistController) {
		v1.POST("/artists/:spotify-id", controller.Add)
		v1.GET("/artists/search", controller.Search)
		subroutes := v1.Group("/artists")
		{
			subroutes.GET("/:artist-id", controller.GetOne)
			registerRoutes(injector, func(controller *apiControllers.ReleaseController) {
				subroutes.GET("/:artist-id/releases/", controller.Get)
				subroutes.GET("/:artist-id/releases/:id/", controller.GetOne)
			})
		}
	})

	registerRoutes(injector, func(controller *apiControllers.StreamingPlatformController) {
		v1.GET("/platforms/sync/collect", controller.ProcessUrlSyncResults)
		v1.GET("/platforms/sync/start", controller.RequestUrlSync)
	})
}

func registerRoutes[T any](injector *do.Injector, function func(*T)) {
	controller := do.MustInvoke[*T](injector)
	function(controller)
}
