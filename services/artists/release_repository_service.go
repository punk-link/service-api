package artists

import (
	artistData "main/data/artists"
	"time"

	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type ReleaseRepositoryService struct {
	artistIdExtractor ArtistIdExtractor
	db                *gorm.DB
	logger            logger.Logger
}

func NewReleaseRepository(injector *do.Injector) (ReleaseRepository, error) {
	artistIdExtractor := do.MustInvoke[ArtistIdExtractor](injector)
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	return &ReleaseRepositoryService{
		artistIdExtractor: artistIdExtractor,
		db:                db,
		logger:            logger,
	}, nil
}

func (t *ReleaseRepositoryService) AddTags(err error, relations *[]artistData.ReleaseTagRelation) error {
	return nil
}

func (t *ReleaseRepositoryService) CreateInBatches(err error, releases *[]artistData.Release) error {
	if err != nil || len(*releases) == 0 {
		return err
	}

	return t.db.Transaction(func(tx *gorm.DB) error {
		err = t.db.CreateInBatches(&releases, CREATE_RELEASES_BATCH_SIZE).Error
		if err != nil {
			return t.handleError(err)
		}

		relations := make([]artistData.ArtistReleaseRelation, 0)
		for _, release := range *releases {
			artistIds := t.artistIdExtractor.ExtractFromOne(release)

			releaseRelations := make([]artistData.ArtistReleaseRelation, 0)
			for _, id := range artistIds {
				releaseRelations = append(releaseRelations, artistData.ArtistReleaseRelation{
					ArtistId:  id,
					ReleaseId: release.Id,
				})
			}

			relations = append(relations, releaseRelations...)
		}

		err = t.db.CreateInBatches(&relations, CREATE_RELATION_BATCH_SIZE).Error
		return t.handleError(err)
	})
}

func (t *ReleaseRepositoryService) Get(err error, artistId int) ([]artistData.Release, error) {
	if err != nil {
		return make([]artistData.Release, 0), err
	}

	subQuery := t.db.Select("release_id").
		Where("artist_id = ?", artistId).
		Table("artist_release_relations")

	var releases []artistData.Release
	err = t.db.Where("id IN (?)", subQuery).
		Order("release_date").
		Find(&releases).
		Error

	return releases, t.handleError(err)
}

func (t *ReleaseRepositoryService) GetCount(err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	var count int64
	t.db.Model(&artistData.Release{}).
		Count(&count)

	return count, t.handleError(err)
}

func (t *ReleaseRepositoryService) GetOne(err error, id int) (artistData.Release, error) {
	if err != nil {
		return artistData.Release{}, err
	}

	var release artistData.Release
	err = t.db.Model(&artistData.Release{}).
		First(&release, id).
		Error

	return release, t.handleError(err)
}

func (t *ReleaseRepositoryService) GetSlimByArtistId(err error, artistId int) ([]artistData.SlimRelease, error) {
	if err != nil {
		return make([]artistData.SlimRelease, 0), err
	}

	subQuery := t.db.Select("release_id").
		Where("artist_id = ?", artistId).
		Table("artist_release_relations")

	var releases []artistData.SlimRelease
	err = t.db.Model(&artistData.Release{}).
		Where("id IN (?)", subQuery).
		Order("release_date").
		Find(&releases).
		Error

	return releases, t.handleError(err)
}

func (t *ReleaseRepositoryService) GetUpcContainers(err error, top int, skip int, updateTreshold time.Time) ([]artistData.Release, error) {
	if err != nil {
		return make([]artistData.Release, 0), err
	}

	var releases []artistData.Release
	err = t.db.Select("id", "upc").
		Where("updated < ?", updateTreshold).
		Order("id").
		Offset(skip).
		Limit(top).
		Find(&releases).
		Error

	return releases, t.handleError(err)
}

func (t *ReleaseRepositoryService) MarksAsUpdated(err error, ids []int, timestamp time.Time) error {
	if err != nil {
		return err
	}

	err = t.db.Model(&artistData.Release{}).
		Where("id in (?)", ids).
		Update("updated", timestamp).
		Error

	return t.handleError(err)
}

func (t *ReleaseRepositoryService) handleError(err error) error {
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return err
}

const CREATE_RELEASES_BATCH_SIZE = 100
const CREATE_RELATION_BATCH_SIZE = 2000
