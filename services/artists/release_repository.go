package artists

import (
	"main/data"
	artistData "main/data/artists"
	"main/services/common"
)

func createDbReleasesInBatches(logger *common.Logger, err error, artists *[]artistData.Release) error {
	if err != nil {
		return err
	}

	err = data.DB.CreateInBatches(&artists, CREATE_RELEASES_BATCH_SIZE).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

func getDbReleasesByArtistId(logger *common.Logger, err error, artistId int) ([]artistData.Release, error) {
	if err != nil {
		return make([]artistData.Release, 0), err
	}

	var releases []artistData.Release
	err = data.DB.Joins("join artist_release_relations rel on rel.release_id = releases.id").
		Where("rel.artist_id = ?", artistId).
		Find(&releases).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return releases, err
}

const CREATE_RELEASES_BATCH_SIZE int = 50
