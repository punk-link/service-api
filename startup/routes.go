package startup

import (
	"main/controllers"
	"main/services/common"
	"main/services/helpers"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func setupRouts(app *gin.Engine, diContainer *dig.Container, logger *common.Logger) {
	err := diContainer.Invoke(func(statusController *controllers.StatusController) {
		app.GET("/health", statusController.CheckHealth)
	})
	if err != nil {
		logDependencyResolvingError(logger, err, controllers.StatusController{})
	}

	v1 := app.Group("/v1")
	err = diContainer.Invoke(func(managerController *controllers.ManagerController) {
		v1.POST("/managers", managerController.AddManager)
		v1.POST("/managers/master", managerController.AddMasterManager)
		v1.GET("/managers", managerController.GetManagers)
		v1.GET("/managers/:id", managerController.GetManager)
		v1.POST("/managers/:id", managerController.ModifyManager)
	})
	if err != nil {
		logDependencyResolvingError(logger, err, controllers.ManagerController{})
	}

	err = diContainer.Invoke(func(labelController *controllers.LabelController) {
		v1.GET("/labels/:id", labelController.GetLabel)
		v1.POST("/labels/:id", labelController.ModifyLabel)
	})
	if err != nil {
		logDependencyResolvingError(logger, err, controllers.LabelController{})
	}

	err = diContainer.Invoke(func(artistController *controllers.ArtistController) {
		v1.GET("/artists/search", artistController.SearchArtist)
		v1.GET("/artists/:spotify-id/releases", artistController.GetReleases)
		v1.GET("/artists/releases/:spotify-id/", artistController.GetRelease)
	})
	if err != nil {
		logDependencyResolvingError(logger, err, controllers.ArtistController{})
	}
}

func logDependencyResolvingError[T any](logger *common.Logger, err error, target T) {
	structName := helpers.GetStructNameAsString(target)
	logger.LogFatal(err, "Can't resolve a dependency '%s': %v", structName, err.Error())
	panic(err.Error())
}
