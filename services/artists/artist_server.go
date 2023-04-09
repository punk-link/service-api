package artists

import (
	artistModels "main/models/artists"
	labelModels "main/models/labels"
)

type ArtistServer interface {
	Add(currentManager labelModels.ManagerContext, spotifyId string) (artistModels.Artist, error)
	Get(labelId int) ([]artistModels.Artist, error)
	GetOne(id int) (artistModels.Artist, error)
	Search(query string) []artistModels.ArtistSearchResult
}
