package artists

import "main/models/common"

type Artist struct {
	Id            int                  `json:"id"`
	ImageMetadata common.ImageMetadata `json:"image"`
	Name          string               `json:"name"`
}
