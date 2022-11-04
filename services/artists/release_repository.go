package artists

import (
	artistData "main/data/artists"
	"time"

	"github.com/punk-link/logger"
	"gorm.io/gorm"
)

func createDbReleasesInBatches(db *gorm.DB, logger logger.Logger, err error, releases *[]artistData.Release) error {
	if err != nil || len(*releases) == 0 {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		err = db.CreateInBatches(&releases, CREATE_RELEASES_BATCH_SIZE).Error
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

		err = db.CreateInBatches(&relations, CREATE_RELATION_BATCH_SIZE).Error
		if err != nil {
			logger.LogError(err, err.Error())
			return err
		}

		return nil
	})
}

func getDbRelease(db *gorm.DB, logger logger.Logger, err error, id int) (artistData.Release, error) {
	if err != nil {
		return artistData.Release{}, err
	}

	var release artistData.Release
	err = db.Model(&artistData.Release{}).
		First(&release, id).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return release, err
}

func getDbReleasesByArtistId(db *gorm.DB, logger logger.Logger, err error, artistId int) ([]artistData.Release, error) {
	if err != nil {
		return make([]artistData.Release, 0), err
	}

	subQuery := db.Select("release_id").
		Where("artist_id = ?", artistId).
		Table("artist_release_relations")

	var releases []artistData.Release
	err = db.Where("id IN (?)", subQuery).
		Order("release_date").
		Find(&releases).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return releases, err
}

func getUpcContainers(db *gorm.DB, logger logger.Logger, err error, top int, skip int, updateTreshold time.Time) ([]artistData.Release, error) {
	if err != nil {
		return make([]artistData.Release, 0), err
	}

	var releases []artistData.Release
	err = db.Select("id", "upc").
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

func getDbReleaseCount(db *gorm.DB, logger logger.Logger, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	var count int64
	db.Model(&artistData.Release{}).
		Count(&count)

	return count, err
}

func markDbReleasesAsUpdated(db *gorm.DB, logger logger.Logger, err error, ids []int, timestamp time.Time) error {
	if err != nil {
		return err
	}

	err = db.Model(&artistData.Release{}).
		Where("id in (?)", ids).
		Update("updated", timestamp).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

const CREATE_RELEASES_BATCH_SIZE = 100
const CREATE_RELATION_BATCH_SIZE = 2000
