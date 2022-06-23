package releases

import (
	"main/models/spotify"
	"main/models/spotify/search"
)

type ArtistRelease struct {
	Id            string                  `json:"id"`
	Artists       []search.Artist         `json:"artists"`
	ImageMetadata []spotify.ImageMetadata `json:"images"`
	Label         string                  `json:"label"`
	Name          string                  `json:"name"`
	ReleaseDate   string                  `json:"release_date"`
	TrackNumber   int                     `json:"total_tracks"`
	Tracks        TrackContainer          `json:"tracks"`
	Type          string                  `json:"album_type"`
}
