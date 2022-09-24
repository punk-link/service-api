package startup

import (
	apiControllers "main/controllers/api"
	mvcControllers "main/controllers/mvc"
	artistServices "main/services/artists"
	"main/services/cache"
	"main/services/common"
	labelServices "main/services/labels"
	spotifyServices "main/services/spotify"

	"github.com/samber/do"
)

func buildDependencies() *do.Injector {
	container := do.New()

	do.Provide(container, common.ConstructLogger)
	do.Provide(container, common.ConstructHashCoder)
	do.Provide(container, cache.ConstructMemoryCacheService)

	do.Provide(container, labelServices.ConstructLabelService)
	do.Provide(container, labelServices.ConstructManagerService)

	do.Provide(container, spotifyServices.ConstructSpotifyService)

	do.Provide(container, artistServices.ConstructReleaseService)
	do.Provide(container, artistServices.ConstructMvcReleaseService)
	do.Provide(container, artistServices.ConstructArtistService)
	do.Provide(container, artistServices.ConstructMvcArtistService)

	do.Provide(container, apiControllers.ConstructArtistController)
	do.Provide(container, apiControllers.ConstructHashController)
	do.Provide(container, apiControllers.ConstructLabelController)
	do.Provide(container, apiControllers.ConstructManagerController)
	do.Provide(container, apiControllers.ConstructReleaseController)
	do.Provide(container, apiControllers.ConstructStatusController)

	do.Provide(container, mvcControllers.ConstructMvcArtistController)
	do.Provide(container, mvcControllers.ConstructMvcReleaseController)

	return container
}
