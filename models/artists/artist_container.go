package artists

import data "main/data/artists"

type ArtistContainer struct {
	Id      string
	Artists []data.Artist
}
