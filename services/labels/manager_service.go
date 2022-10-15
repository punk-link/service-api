package labels

import (
	"main/data"
	labelData "main/data/labels"
	"main/helpers"
	"main/models/labels"

	"main/services/labels/converters"
	"main/services/labels/validators"
	"strings"
	"time"

	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type ManagerService struct {
	labelService *LabelService
	logger       *logger.Logger
}

func ConstructManagerService(injector *do.Injector) (*ManagerService, error) {
	labelService := do.MustInvoke[*LabelService](injector)
	logger := do.MustInvoke[*logger.Logger](injector)

	return &ManagerService{
		labelService: labelService,
		logger:       logger,
	}, nil
}

func (t *ManagerService) Add(currentManager labels.ManagerContext, manager labels.Manager) (labels.Manager, error) {
	trimmedName := strings.TrimSpace(manager.Name)
	err := validators.NameNotEmpty(trimmedName)

	return t.addInternal(err, currentManager, trimmedName)
}

func (t *ManagerService) AddMaster(request labels.AddMasterManagerRequest) (labels.Manager, error) {
	trimmedName := strings.TrimSpace(request.Name)
	err := validators.NameNotEmpty(trimmedName)
	if err != nil {
		return labels.Manager{}, err
	}

	label, err := t.labelService.AddLabel(request.LabelName)
	return t.addInternal(err, labels.ManagerContext{LabelId: label.Id}, trimmedName)
}

func (t *ManagerService) Get(currentManager labels.ManagerContext) ([]labels.Manager, error) {
	dbManagers, err := getDbManagersByLabelId(t.logger, nil, currentManager.LabelId)
	if err != nil {
		return make([]labels.Manager, 0), err
	}

	return converters.ToManagers(dbManagers), nil
}

func (t *ManagerService) GetContext(id int) (labels.ManagerContext, error) {
	manager, err := t.getOneInternal(nil, id)

	return labels.ManagerContext{
		Id:      manager.Id,
		LabelId: manager.LabelId,
	}, err
}

func (t *ManagerService) GetOne(currentManager labels.ManagerContext, id int) (labels.Manager, error) {
	manager, err1 := t.getOneInternal(nil, id)
	err2 := validators.CurrentManagerBelongsToLabel(currentManager, manager.LabelId)

	return manager, helpers.AccumulateErrors(err1, err2)
}

func (t *ManagerService) Modify(currentManager labels.ManagerContext, manager labels.Manager, id int) (labels.Manager, error) {
	err1 := validators.IdConsistsOverRequest(manager.Id, id)

	trimmedName := strings.TrimSpace(manager.Name)
	err2 := validators.NameNotEmpty(trimmedName)

	dbManager, err3 := t.getOneInternal(helpers.AccumulateErrors(err1, err2), manager.Id)
	err4 := validators.CurrentManagerBelongsToLabel(currentManager, dbManager.LabelId)

	return t.modifyInternal(helpers.AccumulateErrors(err3, err4), manager.Id, trimmedName)
}

func (t *ManagerService) addInternal(err error, currentManager labels.ManagerContext, managerName string) (labels.Manager, error) {
	if err != nil {
		return labels.Manager{}, err
	}

	dbManager := converters.ToDbManager(managerName, currentManager.LabelId)
	err = data.DB.Create(&dbManager).Error
	if err != nil {
		t.logger.LogError(err, err.Error())
		return labels.Manager{}, err
	}

	return t.getOneInternal(err, dbManager.Id)
}

func (t *ManagerService) getOneInternal(err error, id int) (labels.Manager, error) {
	if err != nil {
		return labels.Manager{}, err
	}

	var dbManager labelData.Manager
	err = data.DB.First(&dbManager, id).Error
	if err != nil {
		return labels.Manager{}, err
	}

	return converters.ToManager(dbManager), nil
}

func (t *ManagerService) modifyInternal(err error, id int, managerName string) (labels.Manager, error) {
	if err != nil {
		return labels.Manager{}, err
	}

	var dbManager labelData.Manager
	err = data.DB.First(&dbManager, id).Error
	if err != nil {
		return labels.Manager{}, err
	}

	dbManager.Name = managerName
	dbManager.Updated = time.Now().UTC()
	err = data.DB.Save(&dbManager).Error

	return t.getOneInternal(err, id)
}
