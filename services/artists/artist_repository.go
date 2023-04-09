package artists

import artistData "main/data/artists"

type ArtistRepository interface {
	Create(err error, artist *artistData.Artist) error
	CreateInBatches(err error, artists *[]artistData.Artist) error
	Get(err error, ids []int) ([]artistData.Artist, error)
	GetSlim(err error, ids []int) ([]artistData.SlimArtist, error)
	GetBySpotifyIds(err error, spotifyIds []string) ([]artistData.Artist, error)
	GetIdsByLabelId(err error, labelId int) ([]int, error)
	GetOne(err error, id int) (artistData.Artist, error)
	GetOneBySpotifyId(err error, spotifyId string) (artistData.Artist, error)
	GetOneSlim(err error, id int) (artistData.SlimArtist, error)
	Update(err error, artist *artistData.Artist) error
}
