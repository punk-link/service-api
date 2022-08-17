package artists

import "main/models/common"

type Release struct {
	SpotifyId     string                 `json:"spotifyId"`
	Artists       []Artist               `json:"artists"`
	ImageMetadata []common.ImageMetadata `json:"images"`
	Lable         string                 `json:"label"`
	Name          string                 `json:"name"`
	ReleaseDate   string                 `json:"releaseDate"`
	TrackNumber   int                    `json:"trackNumber"`
	Tracks        []Track                `json:"tracks"`
	Type          string                 `json:"type"`
}
