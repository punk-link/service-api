package artists

import (
	"encoding/json"
	artistData "main/data/artists"
	"main/helpers"
	"sort"

	"github.com/punk-link/logger"
)

func getArtistsIdsFromDbRelease(logger logger.Logger, release artistData.Release) []int {
	artistIds := make([]int, 0)

	var featuringArtistIds []int
	featuringArtistErr := json.Unmarshal([]byte(release.FeaturingArtistIds), &featuringArtistIds)
	artistIds = append(artistIds, featuringArtistIds...)

	var releaseArtistIds []int
	releaseArtistErr := json.Unmarshal([]byte(release.FeaturingArtistIds), &releaseArtistIds)
	artistIds = append(artistIds, releaseArtistIds...)

	err := helpers.CombineErrors(featuringArtistErr, releaseArtistErr)
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return artistIds
}

func getArtistsIdsFromDbReleases(logger logger.Logger, releases []artistData.Release) []int {
	artistIds := make([]int, 0)
	for _, release := range releases {
		releaseArtistIds := getArtistsIdsFromDbRelease(logger, release)
		artistIds = append(artistIds, releaseArtistIds...)
	}

	return helpers.Distinct(artistIds)
}

func orderDbReleasesChronologically(target []artistData.Release) {
	sort.Slice(target, func(i, j int) bool {
		return target[i].ReleaseDate.Before(target[j].ReleaseDate)
	})
}
