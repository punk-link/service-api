package labels

import (
	labelData "main/data/labels"
	"time"

	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type LabelRepositoryService struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewLabelRepository(injector *do.Injector) (LabelRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	return &LabelRepositoryService{
		db:     db,
		logger: logger,
	}, nil
}

func (t *LabelRepositoryService) Create(err error, dbLabel *labelData.Label) error {
	if err != nil {
		return err
	}

	err = t.db.Create(&dbLabel).Error
	return t.handleError(err)
}

func (t *LabelRepositoryService) GetOne(err error, id int) (labelData.Label, error) {
	if err != nil {
		return labelData.Label{}, err
	}

	var dbLabel labelData.Label
	err = t.db.First(&dbLabel, id).Error
	return dbLabel, t.handleError(err)
}

func (t *LabelRepositoryService) Update(err error, dbLabel *labelData.Label) error {
	if err != nil {
		return err
	}

	dbLabel.Updated = time.Now().UTC()
	err = t.db.Save(&dbLabel).Error

	return t.handleError(err)
}

func (t *LabelRepositoryService) handleError(err error) error {
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return err
}
