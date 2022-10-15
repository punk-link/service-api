package artists

import (
	"main/data"
	artistData "main/data/artists"
	"time"

	"github.com/punk-link/logger"
	"gorm.io/gorm"
)

func createDbReleasesInBatches(logger *logger.Logger, err error, releases *[]artistData.Release) error {
	if err != nil || len(*releases) == 0 {
		return err
	}

	return data.DB.Transaction(func(tx *gorm.DB) error {
		err = data.DB.CreateInBatches(&releases, CREATE_RELEASES_BATCH_SIZE).Error
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

		err = data.DB.CreateInBatches(&relations, CREATE_RELATION_BATCH_SIZE).Error
		if err != nil {
			logger.LogError(err, err.Error())
			return err
		}

		return nil
	})
}

func getDbRelease(logger *logger.Logger, err error, id int) (artistData.Release, error) {
	if err != nil {
		return artistData.Release{}, err
	}

	var release artistData.Release
	err = data.DB.Model(&artistData.Release{}).
		First(&release, id).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return release, err
}

func getDbReleasesByArtistId(logger *logger.Logger, err error, artistId int) ([]artistData.Release, error) {
	if err != nil {
		return make([]artistData.Release, 0), err
	}

	subQuery := data.DB.Select("release_id").
		Where("artist_id = ?", artistId).
		Table("artist_release_relations")

	var releases []artistData.Release
	err = data.DB.Where("id IN (?)", subQuery).
		Order("release_date").
		Find(&releases).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return releases, err
}

func getUpcContainers(logger *logger.Logger, err error, top int, skip int, updateTreshold time.Time) ([]artistData.Release, error) {
	if err != nil {
		return make([]artistData.Release, 0), err
	}

	var releases []artistData.Release
	err = data.DB.Select("id", "upc").
		Where("updated < ?", updateTreshold).
		Order("id").
		Offset(skip).
		Limit(top).
		Find(&releases).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return releases, err
}

func getDbReleaseCount(logger *logger.Logger, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	var count int64
	data.DB.Model(&artistData.Release{}).
		Count(&count)

	return count, err
}

func markDbReleasesAsUpdated(logger *logger.Logger, err error, ids []int, timestamp time.Time) error {
	if err != nil {
		return err
	}

	err = data.DB.Model(&artistData.Release{}).
		Where("id in (?)", ids).
		Update("updated", timestamp).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

const CREATE_RELEASES_BATCH_SIZE int = 100
const CREATE_RELATION_BATCH_SIZE int = 2000
