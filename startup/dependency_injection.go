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

	container.Provide(common.ConstructLogger)

	container.Provide(labelServices.ConstructLabelService)
	container.Provide(labelServices.ConstructManagerService)

	container.Provide(spotifyServices.ConstructSpotifyService)

	container.Provide(artistServices.ConstructReleaseService)
	container.Provide(artistServices.ConstructArtistCacheService)
	container.Provide(artistServices.ConstructArtistService)

	container.Provide(controllers.ConstructArtistController)
	container.Provide(controllers.ConstructReleaseController)
	container.Provide(controllers.ConstructLabelController)
	container.Provide(controllers.ConstructManagerController)
	container.Provide(controllers.ConstructStatusController)

	return container
}
