package labels

import (
	"main/data"
	labelData "main/data/labels"
	"main/models/labels"
	requests "main/requests/labels"
	"main/services/helpers"
	validator "main/services/labels/validators"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func AddManager(manager labels.Manager, currentManager labels.ManagerContext) (labels.Manager, error) {
	trimmedName := strings.TrimSpace(manager.Name)
	if err := validator.NameNotEmpty(trimmedName); err != nil {
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

	return GetManager(dbManager.Id, currentManager)
}

func AddMasterManager(request requests.AddMasterManagerRequest) (labels.Manager, error) {
	label, err := AddLabel(request.LabelName)
	if err != nil {
		return labels.Manager{}, err
	}

	manager, err := AddManager(labels.Manager{Name: request.Name}, labels.ManagerContext{LabelId: label.Id})
	if err != nil {
		return labels.Manager{}, err
	}

	currentManager, _ := GetManagerContext(manager.Id) // Assuming there is no error here
	return GetManager(manager.Id, currentManager)
}

func GetManager(id int, currentManager labels.ManagerContext) (labels.Manager, error) {
	dbManager, err := helpers.GetData[labelData.Manager](id)
	if err != nil {
		return labels.Manager{}, err
	}

	if err := validator.CurrentManagerBelongsToLabel(currentManager, dbManager.LabelId); err != nil {
		return labels.Manager{}, err
	}

	return toManager(dbManager), nil
}

func GetManagerContext(id int) (labels.ManagerContext, error) {
	dbManager, err := helpers.GetData[labelData.Manager](id)
	if err != nil {
		return labels.ManagerContext{}, err
	}

	return labels.ManagerContext{
		Id:      dbManager.Id,
		LabelId: dbManager.LabelId,
	}, nil
}

func GetLabelManagers(currentManager labels.ManagerContext) []labels.Manager {
	var dbManagers []labelData.Manager
	data.DB.Where("label_id = ?", currentManager.LabelId).Find(&dbManagers)

	managers := make([]labels.Manager, len(dbManagers))
	for i, manager := range dbManagers {
		managers[i] = toManager(manager)
	}

	return managers
}

func ModifyManager(manager labels.Manager, id int, currentManager labels.ManagerContext) (labels.Manager, error) {
	if err := validator.IdConsistsOverRequest(manager.Id, id); err != nil {
		return labels.Manager{}, err
	}

	trimmedName := strings.TrimSpace(manager.Name)
	if err := validator.NameNotEmpty(trimmedName); err != nil {
		return labels.Manager{}, err
	}

	dbManager, err := helpers.GetData[labelData.Manager](id)
	if err != nil {
		return labels.Manager{}, err
	}

	if err := validator.CurrentManagerBelongsToLabel(currentManager, dbManager.LabelId); err != nil {
		return labels.Manager{}, err
	}

	dbManager.Name = trimmedName
	dbManager.Updated = time.Now().UTC()
	data.DB.Save(&dbManager)

	return GetManager(id, currentManager)
}

func toManager(source labelData.Manager) labels.Manager {
	return labels.Manager{
		Id:      source.Id,
		Name:    source.Name,
		LabelId: source.LabelId,
	}
}
