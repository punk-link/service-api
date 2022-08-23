package artists

import "main/models/common"

type Artist struct {
	Id           int                 `json:"id"`
	ImageDetails common.ImageDetails `json:"image"`
	LabelId      int                 `json:"labelId"`
	Name         string              `json:"name"`
	Releases     []Release           `json:"releases"`
}
