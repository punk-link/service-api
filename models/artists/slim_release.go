package artists

import (
	"main/models/common"
	"time"
)

type SlimRelease struct {
	Id           int                 `json:"id"`
	ImageDetails common.ImageDetails `json:"image"`
	Name         string              `json:"name"`
	ReleaseDate  time.Time           `json:"releaseDate"`
	Type         string              `json:"type"`
}
