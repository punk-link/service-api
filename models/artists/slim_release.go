package artists

import (
	"main/models/common"
	"time"
)

type SlimRelease struct {
	Slug         string              `json:"slug"`
	ImageDetails common.ImageDetails `json:"image"`
	Name         string              `json:"name"`
	ReleaseDate  time.Time           `json:"releaseDate"`
	Type         string              `json:"type"`
}
