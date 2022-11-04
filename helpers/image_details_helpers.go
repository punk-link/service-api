package helpers

import (
	"main/models/platforms/spotify"
	"sort"
)

func OrderImageDetailsDesc(target []spotify.ImageDetails) {
	sort.SliceStable(target, func(i, j int) bool {
		return target[i].Height > target[j].Height
	})
}
