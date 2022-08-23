package helpers

import (
	"main/models/spotify"
	"sort"
)

func ReorderImageDetailsDesc(target []spotify.ImageDetails) {
	sort.SliceStable(target, func(i, j int) bool {
		return target[i].Height > target[j].Height
	})
}
