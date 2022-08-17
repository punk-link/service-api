package search

import "main/models/spotify"

type Artist struct {
	Id           string                 `json:"id"`
	ImageDetails []spotify.ImageDetails `json:"images"`
	Name         string                 `json:"name"`
}
