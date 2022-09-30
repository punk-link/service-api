package releases

import (
	"main/models/spotify"
	"main/models/spotify/artists"
)

type Release struct {
	Id                   string                 `json:"id"`
	Artists              []artists.SlimArtist   `json:"artists"`
	ExternalIds          ExternalIds            `json:"external_ids"`
	ExternalUrls         ExternalUrls           `json:"external_urls"`
	ImageDetails         []spotify.ImageDetails `json:"images"`
	Label                string                 `json:"label"`
	Name                 string                 `json:"name"`
	ReleaseDate          string                 `json:"release_date"`
	ReleaseDatePrecision string                 `json:"release_date_precision"`
	TrackNumber          int                    `json:"total_tracks"`
	Tracks               TrackContainer         `json:"tracks"`
	Type                 string                 `json:"album_type"`
}
