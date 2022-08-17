package releases

import (
	"main/models/spotify"
	"main/models/spotify/search"
)

type Release struct {
	Id                   string                 `json:"id"`
	Artists              []search.Artist        `json:"artists"`
	ImageDetails         []spotify.ImageDetails `json:"images"`
	Label                string                 `json:"label"`
	Name                 string                 `json:"name"`
	ReleaseDate          string                 `json:"release_date"`
	ReleaseDatePrecision string                 `json:"release_date_precision"`
	TrackNumber          int                    `json:"total_tracks"`
	Tracks               TrackContainer         `json:"tracks"`
	Type                 string                 `json:"album_type"`
}
