package search

import "main/models/spotify"

type Artist struct {
	Id            string                  `json:"id"`
	ImageMetadata []spotify.ImageMetadata `json:"images"`
	Name          string                  `json:"name"`
}
