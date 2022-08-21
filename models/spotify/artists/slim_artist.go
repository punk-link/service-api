package artists

import "main/models/spotify"

type SlimArtist struct {
	Id           string                 `json:"id"`
	ImageDetails []spotify.ImageDetails `json:"images"`
	Name         string                 `json:"name"`
}
