package startup

import (
	apiControllers "main/controllers/api"
	mvcControllers "main/controllers/mvc"
	artistServices "main/services/artists"
	"main/services/cache"
	"main/services/common"
	labelServices "main/services/labels"
	spotifyServices "main/services/spotify"

	"go.uber.org/dig"
)

func buildDependencies() *dig.Container {
	container := dig.New()

	container.Provide(common.ConstructLogger)
	container.Provide(cache.ConstructMemoryCacheService)

	container.Provide(labelServices.ConstructLabelService)
	container.Provide(labelServices.ConstructManagerService)

	container.Provide(spotifyServices.ConstructSpotifyService)

	container.Provide(artistServices.ConstructReleaseService)
	container.Provide(artistServices.ConstructMvcReleaseService)
	container.Provide(artistServices.ConstructArtistService)

	container.Provide(apiControllers.ConstructArtistController)
	container.Provide(apiControllers.ConstructLabelController)
	container.Provide(apiControllers.ConstructManagerController)
	container.Provide(apiControllers.ConstructReleaseController)
	container.Provide(apiControllers.ConstructStatusController)

	container.Provide(mvcControllers.ConstructMvcReleaseController)

	return container
}
