package labels

import (
	"main/data"
	labelData "main/data/labels"
	"main/services/common"
)

func getDbManagersByLabelId(logger *common.Logger, err error, labelId int) ([]labelData.Manager, error) {
	if err != nil {
		return make([]labelData.Manager, 0), err
	}

	var dbManagers []labelData.Manager
	err = data.DB.Where("label_id = ?", labelId).
		Find(&dbManagers).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return dbManagers, err
}
