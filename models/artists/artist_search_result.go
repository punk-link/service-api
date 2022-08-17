package artists

import "main/models/common"

type ArtistSearchResult struct {
	SpotifyId    string              `json:"spotifyId"`
	ImageDetails common.ImageDetails `json:"image"`
	Name         string              `json:"name"`
}
