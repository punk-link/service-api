package converters

import (
	"main/helpers"
	"main/models/common"
	"main/models/spotify"
)

func ToImageDetailsFromSpotify(imageDetails []spotify.ImageDetails) common.ImageDetails {
	if len(imageDetails) == 0 {
		return common.ImageDetails{}
	}

	helpers.ReorderImageDetailsDesc(imageDetails)
	return common.ImageDetails{
		Height: imageDetails[0].Height,
		Url:    imageDetails[0].Url,
	}
}
