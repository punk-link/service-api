package startup

import (
	"main/controllers"
	artistServices "main/services/artists"
	"main/services/common"
	labelServices "main/services/labels"
	spotifyServices "main/services/spotify"

	"go.uber.org/dig"
)

func buildDependencies() *dig.Container {
	container := dig.New()

	container.Provide(common.BuildLogger)
	container.Provide(labelServices.BuildLabelService)
	container.Provide(labelServices.BuildManagerService)
	container.Provide(spotifyServices.BuildSpotifyService)
	container.Provide(artistServices.BuildArtistService)

	container.Provide(controllers.BuildArtistController)
	container.Provide(controllers.BuildLabelController)
	container.Provide(controllers.BuildManagerController)
	container.Provide(controllers.BuildStatusController)

	return container
}
