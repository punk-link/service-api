package converters

import (
	"main/models/common"
	"main/models/spotify"
	"sort"
)

func ToImageDetailsFromSpotify(imageDetails []spotify.ImageDetails) common.ImageDetails {
	if len(imageDetails) == 0 {
		return common.ImageDetails{}
	}

	sort.SliceStable(imageDetails, func(i, j int) bool {
		return imageDetails[i].Height > imageDetails[j].Height
	})

	return common.ImageDetails{
		Height: imageDetails[0].Height,
		Url:    imageDetails[0].Url,
	}
}
