package labels

import (
	"main/data"
	labelData "main/data/labels"
	"main/models/labels"
	requests "main/requests/labels"
	"main/services/helpers"
	"time"

	"github.com/rs/zerolog/log"
)

func AddManager(currentManager labels.Manager, manager labels.Manager) (labels.Manager, error) {
	trimmedName, err := validateAndTrimName(manager.Name)
	if err != nil {
		return labels.Manager{}, err
	}

	now := time.Now().UTC()
	dbManager := labelData.Manager{
		Created: now,
		Name:    trimmedName,
		LabelId: currentManager.LabelId,
		Updated: now,
	}

	result := data.DB.Create(&dbManager)
	if result.Error != nil {
		log.Logger.Error().Err(result.Error).Msg(result.Error.Error())
		return labels.Manager{}, result.Error
	}

	return GetManager(dbManager.Id)
}

func AddMasterManager(request requests.AddMasterManagerRequest) (labels.Manager, error) {
	label, err := AddLabel(request.LabelName)
	if err != nil {
		return labels.Manager{}, err
	}

	manager, err := AddManager(labels.Manager{LabelId: label.Id}, labels.Manager{Name: request.Name})
	if err != nil {
		return labels.Manager{}, err
	}

	return GetManager(manager.Id)
}

func GetManager(id int) (labels.Manager, error) {
	dbManager, err := helpers.GetData[labelData.Manager](id)
	if err != nil {
		return labels.Manager{}, err
	}

	return toManager(dbManager), nil
}

func GetLabelManagers(labelId int) []labels.Manager {
	var dbManagers []labelData.Manager
	data.DB.Where("label_id = ?", labelId).Find(&dbManagers)

	managers := make([]labels.Manager, len(dbManagers))
	for i, manager := range dbManagers {
		managers[i] = toManager(manager)
	}

	return managers
}

func ModifyManager(manager labels.Manager, id int) (labels.Manager, error) {
	if err := validateId(manager.Id, id); err != nil {
		return labels.Manager{}, err
	}

	trimmedName, err := validateAndTrimName(manager.Name)
	if err != nil {
		return labels.Manager{}, err
	}

	dbManager, err := helpers.GetData[labelData.Manager](id)
	if err != nil {
		return labels.Manager{}, err
	}

	dbManager.Name = trimmedName
	dbManager.Updated = time.Now().UTC()
	data.DB.Save(&dbManager)

	return GetManager(id)
}

func toManager(source labelData.Manager) labels.Manager {
	return labels.Manager{
		Id:      source.Id,
		Name:    source.Name,
		LabelId: source.LabelId,
	}
}
