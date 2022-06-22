package releases

import "main/models/spotify"

type ArtistRelease struct {
	Id            string                  `json:"id"`
	ImageMetadata []spotify.ImageMetadata `json:"images"`
	Name          string
	ReleaseDate   string `json:"release_date"`
	TrackNumber   int    `json:"total_tracks"`
	Type          string
}
