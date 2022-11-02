package searches

import "main/models/platforms/spotify/artists"

type ArtistContainer struct {
	Items []artists.SlimArtist `json:"items"`
}
