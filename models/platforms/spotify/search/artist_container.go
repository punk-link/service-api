package search

import "main/models/platforms/spotify/artists"

type ArtistContainer struct {
	Items []artists.SlimArtist `json:"items"`
}
