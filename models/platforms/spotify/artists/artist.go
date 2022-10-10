package artists

import "main/models/platforms/spotify"

type Artist struct {
	Id           string                 `json:"id"`
	Genres       []string               `json:"genres"`
	ImageDetails []spotify.ImageDetails `json:"images"`
	Name         string                 `json:"name"`
}
