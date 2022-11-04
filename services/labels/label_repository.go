package labels

import (
	labelData "main/data/labels"
	"time"

	"github.com/punk-link/logger"
	"gorm.io/gorm"
)

func createDbLabel(db *gorm.DB, logger logger.Logger, err error, dbLabel *labelData.Label) error {
	if err != nil {
		return err
	}

	err = db.Create(&dbLabel).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

func getDbLabel(db *gorm.DB, logger logger.Logger, err error, id int) (labelData.Label, error) {
	if err != nil {
		return labelData.Label{}, err
	}

	var dbLabel labelData.Label
	err = db.First(&dbLabel, id).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return dbLabel, err
}

func getDbManagersByLabelId(db *gorm.DB, logger logger.Logger, err error, labelId int) ([]labelData.Manager, error) {
	if err != nil {
		return make([]labelData.Manager, 0), err
	}

	var dbManagers []labelData.Manager
	err = db.Where("label_id = ?", labelId).
		Find(&dbManagers).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return dbManagers, err
}

func updateDbLabel(db *gorm.DB, logger logger.Logger, err error, dbLabel *labelData.Label) error {
	if err != nil {
		return err
	}

	dbLabel.Updated = time.Now().UTC()
	err = db.Save(&dbLabel).Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}
