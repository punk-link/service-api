package artists

import "main/models/common"

type ArtistSearchResult struct {
	SpotifyId     string               `json:"spotifyId"`
	ImageMetadata common.ImageMetadata `json:"image"`
	Name          string               `json:"name"`
}
