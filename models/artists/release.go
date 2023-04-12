package artists

import (
	"main/models/common"
	"time"
)

type Release struct {
	Id               int                 `json:"id"`
	Description      string              `json:"description"`
	FeaturingArtists []Artist            `json:"featuringArtists"`
	ImageDetails     common.ImageDetails `json:"image"`
	Label            string              `json:"label"`
	Name             string              `json:"name"`
	ReleaseArtists   []Artist            `json:"releaseArtists"`
	ReleaseDate      time.Time           `json:"releaseDate"`
	TrackNumber      int                 `json:"trackNumber"`
	Tracks           []Track             `json:"tracks"`
	Type             string              `json:"type"`
}
