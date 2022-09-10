package artists

import (
	"main/data"
	artistData "main/data/artists"
	"main/services/common"

	"gorm.io/gorm"
)

func createDbReleasesInBatches(logger *common.Logger, err error, releases *[]artistData.Release) error {
	if err != nil || len(*releases) == 0 {
		return err
	}

	return data.DB.Transaction(func(tx *gorm.DB) error {
		err = data.DB.CreateInBatches(&releases, CREATE_RELEASES_BATCH_SIZE).
			Error

		if err != nil {
			logger.LogError(err, err.Error())
			return err
		}

		relations := make([]artistData.ArtistReleaseRelation, 0)
		for _, release := range *releases {
			artistIds := getArtistsIdsFromDbRelease(logger, release)

			releaseRelations := make([]artistData.ArtistReleaseRelation, 0)
			for _, id := range artistIds {
				releaseRelations = append(releaseRelations, artistData.ArtistReleaseRelation{
					ArtistId:  id,
					ReleaseId: release.Id,
				})
			}

			relations = append(relations, releaseRelations...)
		}

		err = data.DB.CreateInBatches(&relations, CREATE_RELEASES_BATCH_SIZE).
			Error

		if err != nil {
			logger.LogError(err, err.Error())
			return err
		}

		return nil
	})
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
