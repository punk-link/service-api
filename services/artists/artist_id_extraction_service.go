package artists

import (
	"encoding/json"
	artistData "main/data/artists"
	"main/helpers"

	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type ArtistIdExtractingService struct {
	logger logger.Logger
}

func NewArtistIdExtractor(injector *do.Injector) (ArtistIdExtractor, error) {
	logger := do.MustInvoke[logger.Logger](injector)

	return &ArtistIdExtractingService{
		logger: logger,
	}, nil
}

func (t *ArtistIdExtractingService) Extract(releases []artistData.Release) []int {
	artistIds := make([]int, 0)
	for _, release := range releases {
		releaseArtistIds := t.ExtractFromOne(release)
		artistIds = append(artistIds, releaseArtistIds...)
	}

	return helpers.Distinct(artistIds)
}

func (t *ArtistIdExtractingService) ExtractFromOne(release artistData.Release) []int {
	artistIds := make([]int, 0)

	var featuringArtistIds []int
	featuringArtistErr := json.Unmarshal([]byte(release.FeaturingArtistIds), &featuringArtistIds)
	artistIds = append(artistIds, featuringArtistIds...)

	var releaseArtistIds []int
	releaseArtistErr := json.Unmarshal([]byte(release.FeaturingArtistIds), &releaseArtistIds)
	artistIds = append(artistIds, releaseArtistIds...)

	err := helpers.CombineErrors(featuringArtistErr, releaseArtistErr)
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return artistIds
}
