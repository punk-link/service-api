package repositories

import (
	platformData "main/data/platforms"
	"main/helpers"

	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type PlatformUrlRepositoryService struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewPlatformUrlRepository(injector *do.Injector) (PlatformUrlRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	return &PlatformUrlRepositoryService{
		db:     db,
		logger: logger,
	}, nil
}

func (t *PlatformUrlRepositoryService) CreatelsInBatches(err error, urls []platformData.PlatformReleaseUrl) error {
	if err != nil {
		return err
	}

	err = t.db.CreateInBatches(&urls, CREATE_PLATFORM_RELEASE_URLS_BATCH_SIZE).Error
	return t.handleError(err)
}

func (t *PlatformUrlRepositoryService) GetByReleaseId(err error, id int) ([]platformData.PlatformReleaseUrl, error) {
	return t.GetByReleaseIds(err, []int{id})
}

func (t *PlatformUrlRepositoryService) GetByReleaseIds(err error, ids []int) ([]platformData.PlatformReleaseUrl, error) {
	if err != nil {
		return make([]platformData.PlatformReleaseUrl, 0), err
	}

	var results []platformData.PlatformReleaseUrl
	err = t.db.Where("release_id in (?)", ids).
		Find(&results).
		Error

	return results, t.handleError(err)
}

func (t *PlatformUrlRepositoryService) UpdateInBatches(err error, urls []platformData.PlatformReleaseUrl) error {
	if err != nil || len(urls) == 0 {
		return err
	}

	for _, url := range urls {
		innerErr := t.db.Model(&platformData.PlatformReleaseUrl{}).
			Where("id = ?", url.Id).
			Updates(url).
			Error

		err = helpers.CombineErrors(err, innerErr)
	}

	return t.handleError(err)
}

func (t *PlatformUrlRepositoryService) handleError(err error) error {
	if helpers.ShouldHandleDbError(err) {
		t.logger.LogError(err, err.Error())
	}

	return err
}

const CREATE_PLATFORM_RELEASE_URLS_BATCH_SIZE = 200
