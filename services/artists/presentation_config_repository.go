package artists

import artistData "main/data/artists"

type PresentationConfigRepository interface {
	Get(err error, artistId int) (artistData.ArtistPresentationConfig, error)
}
