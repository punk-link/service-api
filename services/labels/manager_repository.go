package labels

import (
	labelData "main/data/labels"
	"time"

	"github.com/punk-link/logger"
	"gorm.io/gorm"
)

func createDbManager(db *gorm.DB, logger logger.Logger, err error, dbManager *labelData.Manager) error {
	if err != nil {
		return err
	}

	err = db.Create(dbManager).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

func getDbManager(db *gorm.DB, logger logger.Logger, err error, id int) (labelData.Manager, error) {
	if err != nil {
		return labelData.Manager{}, err
	}

	var dbManager labelData.Manager
	err = db.First(&dbManager, id).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return dbManager, err
}

func updateDbManager(db *gorm.DB, logger logger.Logger, err error, dbManager *labelData.Manager) error {
	if err != nil {
		return err
	}

	dbManager.Updated = time.Now().UTC()
	err = db.Save(&dbManager).Error

	return err
}
