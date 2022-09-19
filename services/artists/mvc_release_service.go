package artists

import (
	"main/services/cache"
	"main/services/common"
)

type MvcReleaseService struct {
	cache  *cache.MemoryCacheService
	logger *common.Logger
}

func ConstructMvcReleaseService(cache *cache.MemoryCacheService, logger *common.Logger) *MvcReleaseService {
	return &MvcReleaseService{
		cache:  cache,
		logger: logger,
	}
}

func (t *MvcReleaseService) Get(hash string) (map[string]any, error) {
	//templ := template.New("Release")
	//x := templ.Execute()

	return map[string]any{
		"PageTitle":    "The Indeepandas - Planescape of Limits",
		"ArtistName":   "The Indeepandas",
		"ReleaseTitle": "Planescape of Limits",
		"ReleaseDate":  2021,
		"ImageUrl":     "https://www.chertanovo.link/img/the-indeepandas/the-indeepandas-planescape-of-limits-large.jpg",
		"Tracks":       []string{"Varg E", "На Думской E", "Animal Farm E", "Welcome Home E", "37", "At Home Among Strangers", "Tolerate Intolerance"},
		"Services":     []string{"Apple Music", "Deezer"},
	}, nil
}
