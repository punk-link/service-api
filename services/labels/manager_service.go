package labels

import (
	"main/helpers"
	labelModels "main/models/labels"

	"main/services/labels/converters"
	"main/services/labels/repositories"
	"main/services/labels/validators"
	"strings"

	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type ManagerService struct {
	labelService LabelServer
	logger       logger.Logger
	repository   repositories.ManagerRepository
}

func NewManagerService(injector *do.Injector) (ManagerServer, error) {
	labelService := do.MustInvoke[LabelServer](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	repository := do.MustInvoke[repositories.ManagerRepository](injector)

	return &ManagerService{
		labelService: labelService,
		logger:       logger,
		repository:   repository,
	}, nil
}

func (t *ManagerService) Add(currentManager labelModels.ManagerContext, manager labelModels.Manager) (labelModels.Manager, error) {
	trimmedName := strings.TrimSpace(manager.Name)
	err := validators.NameNotEmpty(trimmedName)

	return t.addInternal(err, currentManager, trimmedName)
}

func (t *ManagerService) AddMaster(request labelModels.AddMasterManagerRequest) (labelModels.Manager, error) {
	trimmedName := strings.TrimSpace(request.Name)
	err := validators.NameNotEmpty(trimmedName)
	if err != nil {
		return labelModels.Manager{}, err
	}

	label, err := t.labelService.Add(request.LabelName)
	return t.addInternal(err, labelModels.ManagerContext{LabelId: label.Id}, trimmedName)
}

func (t *ManagerService) Get(currentManager labelModels.ManagerContext) ([]labelModels.Manager, error) {
	dbManagers, err := t.repository.GetByLabelId(nil, currentManager.LabelId)
	if err != nil {
		return make([]labelModels.Manager, 0), err
	}

	return converters.ToManagers(dbManagers), nil
}

func (t *ManagerService) GetContext(id int) (labelModels.ManagerContext, error) {
	manager, err := t.getOneInternal(nil, id)

	return labelModels.ManagerContext{
		Id:      manager.Id,
		LabelId: manager.LabelId,
	}, err
}

func (t *ManagerService) GetOne(currentManager labelModels.ManagerContext, id int) (labelModels.Manager, error) {
	manager, err1 := t.getOneInternal(nil, id)
	err2 := validators.CurrentManagerBelongsToLabel(currentManager, manager.LabelId)

	return manager, helpers.AccumulateErrors(err1, err2)
}

func (t *ManagerService) Modify(currentManager labelModels.ManagerContext, manager labelModels.Manager, id int) (labelModels.Manager, error) {
	err1 := validators.IdConsistsOverRequest(manager.Id, id)

	trimmedName := strings.TrimSpace(manager.Name)
	err2 := validators.NameNotEmpty(trimmedName)

	dbManager, err3 := t.getOneInternal(helpers.AccumulateErrors(err1, err2), manager.Id)
	err4 := validators.CurrentManagerBelongsToLabel(currentManager, dbManager.LabelId)

	return t.modifyInternal(helpers.AccumulateErrors(err3, err4), manager.Id, trimmedName)
}

func (t *ManagerService) addInternal(err error, currentManager labelModels.ManagerContext, managerName string) (labelModels.Manager, error) {
	if err != nil {
		return labelModels.Manager{}, err
	}

	dbManager := converters.ToDbManager(managerName, currentManager.LabelId)
	err = t.repository.Create(err, &dbManager)
	return t.getOneInternal(err, dbManager.Id)
}

func (t *ManagerService) getOneInternal(err error, id int) (labelModels.Manager, error) {
	if err != nil {
		return labelModels.Manager{}, err
	}

	dbManager, err := t.repository.GetOne(err, id)
	return converters.ToManager(err, dbManager)
}

func (t *ManagerService) modifyInternal(err error, id int, managerName string) (labelModels.Manager, error) {
	if err != nil {
		return labelModels.Manager{}, err
	}

	dbManager, err := t.repository.GetOne(err, id)

	dbManager.Name = managerName
	err = t.repository.Update(err, &dbManager)
	return t.getOneInternal(err, id)
}
