package repositories

import (
	labelData "main/data/labels"
	"time"

	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type ManagerRepositoryService struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewManagerRepository(injector *do.Injector) (ManagerRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	return &ManagerRepositoryService{
		db:     db,
		logger: logger,
	}, nil
}

func (t *ManagerRepositoryService) Create(err error, dbManager *labelData.Manager) error {
	if err != nil {
		return err
	}

	err = t.db.Create(dbManager).Error
	return t.handleError(err)
}

func (t *ManagerRepositoryService) GetByLabelId(err error, labelId int) ([]labelData.Manager, error) {
	if err != nil {
		return make([]labelData.Manager, 0), err
	}

	var dbManagers []labelData.Manager
	err = t.db.Where("label_id = ?", labelId).
		Find(&dbManagers).
		Error

	return dbManagers, t.handleError(err)
}

func (t *ManagerRepositoryService) GetOne(err error, id int) (labelData.Manager, error) {
	if err != nil {
		return labelData.Manager{}, err
	}

	var dbManager labelData.Manager
	err = t.db.First(&dbManager, id).Error
	return dbManager, t.handleError(err)
}

func (t *ManagerRepositoryService) Update(err error, dbManager *labelData.Manager) error {
	if err != nil {
		return err
	}

	dbManager.Updated = time.Now().UTC()
	err = t.db.Save(&dbManager).Error

	return t.handleError(err)
}

func (t *ManagerRepositoryService) handleError(err error) error {
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return err
}
