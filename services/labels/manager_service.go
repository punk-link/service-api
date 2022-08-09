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

type ManagerService struct {
	labelService *LabelService
}

func NewManagerService(labelService *LabelService) *ManagerService {
	return &ManagerService{
		labelService: labelService,
	}
}

func (service *ManagerService) AddManager(currentManager labels.ManagerContext, manager labels.Manager) (labels.Manager, error) {
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

	return service.GetManager(currentManager, dbManager.Id)
}

func (service *ManagerService) AddMasterManager(request requests.AddMasterManagerRequest) (labels.Manager, error) {
	label, err := service.labelService.AddLabel(request.LabelName)
	if err != nil {
		return labels.Manager{}, err
	}

	manager, err := service.AddManager(labels.ManagerContext{LabelId: label.Id}, labels.Manager{Name: request.Name})
	if err != nil {
		return labels.Manager{}, err
	}

	currentManager, _ := service.GetManagerContext(manager.Id) // Assuming there is no error here
	return service.GetManager(currentManager, manager.Id)
}

func (service *ManagerService) GetManager(currentManager labels.ManagerContext, id int) (labels.Manager, error) {
	dbManager, err := helpers.GetEntity[labelData.Manager](id)
	if err != nil {
		return labels.Manager{}, err
	}

	if err := validator.CurrentManagerBelongsToLabel(currentManager, dbManager.LabelId); err != nil {
		return labels.Manager{}, err
	}

	return toManager(dbManager), nil
}

func (service *ManagerService) GetManagerContext(id int) (labels.ManagerContext, error) {
	dbManager, err := helpers.GetEntity[labelData.Manager](id)
	if err != nil {
		return labels.ManagerContext{}, err
	}

	return labels.ManagerContext{
		Id:      dbManager.Id,
		LabelId: dbManager.LabelId,
	}, nil
}

func (service *ManagerService) GetLabelManagers(currentManager labels.ManagerContext) []labels.Manager {
	var dbManagers []labelData.Manager
	data.DB.Where("label_id = ?", currentManager.LabelId).Find(&dbManagers)

	managers := make([]labels.Manager, len(dbManagers))
	for i, manager := range dbManagers {
		managers[i] = toManager(manager)
	}

	return managers
}

func (service *ManagerService) ModifyManager(currentManager labels.ManagerContext, manager labels.Manager, id int) (labels.Manager, error) {
	if err := validator.IdConsistsOverRequest(manager.Id, id); err != nil {
		return labels.Manager{}, err
	}

	trimmedName := strings.TrimSpace(manager.Name)
	if err := validator.NameNotEmpty(trimmedName); err != nil {
		return labels.Manager{}, err
	}

	dbManager, err := helpers.GetEntity[labelData.Manager](id)
	if err != nil {
		return labels.Manager{}, err
	}

	if err := validator.CurrentManagerBelongsToLabel(currentManager, dbManager.LabelId); err != nil {
		return labels.Manager{}, err
	}

	dbManager.Name = trimmedName
	dbManager.Updated = time.Now().UTC()
	data.DB.Save(&dbManager)

	return service.GetManager(currentManager, id)
}

func toManager(source labelData.Manager) labels.Manager {
	return labels.Manager{
		Id:      source.Id,
		Name:    source.Name,
		LabelId: source.LabelId,
	}
}
