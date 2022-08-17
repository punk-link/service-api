package artists

import "main/models/common"

type Artist struct {
	Id           int                 `json:"id"`
	ImageDetails common.ImageDetails `json:"image"`
	Name         string              `json:"name"`
}
