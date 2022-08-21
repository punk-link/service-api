package search

import "main/models/spotify/artists"

type ArtistContainer struct {
	Items []artists.SlimArtist `json:"items"`
}
