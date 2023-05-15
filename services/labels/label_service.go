package labels

import (
	labelData "main/data/labels"
	"main/helpers"
	labelModels "main/models/labels"
	"main/services/labels/repositories"
	"main/services/labels/validators"
	"strings"
	"time"

	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type LabelService struct {
	logger     logger.Logger
	repository repositories.LabelRepository
}

func NewLabelService(injector *do.Injector) (LabelServer, error) {
	logger := do.MustInvoke[logger.Logger](injector)
	repository := do.MustInvoke[repositories.LabelRepository](injector)

	return &LabelService{
		logger:     logger,
		repository: repository,
	}, nil
}

func (t *LabelService) Add(labelName string) (labelModels.Label, error) {
	trimmedName := strings.TrimSpace(labelName)
	err := validators.NameNotEmpty(trimmedName)

	return t.addInternal(err, trimmedName)
}

func (t *LabelService) GetOne(currentManager labelModels.ManagerContext, id int) (labelModels.Label, error) {
	err := validators.CurrentManagerBelongsToLabel(currentManager, id)
	return t.getWithoutContextCheck(err, id)
}

func (t *LabelService) Modify(currentManager labelModels.ManagerContext, label labelModels.Label, id int) (labelModels.Label, error) {
	err1 := validators.CurrentManagerBelongsToLabel(currentManager, id)
	err2 := validators.IdConsistsOverRequest(label.Id, id)

	trimmedName := strings.TrimSpace(label.Name)
	err3 := validators.NameNotEmpty(trimmedName)

	return t.modifyInternal(helpers.AccumulateErrors(err1, err2, err3), currentManager, trimmedName)
}

func (t *LabelService) addInternal(err error, labelName string) (labelModels.Label, error) {
	if err != nil {
		return labelModels.Label{}, err
	}

	now := time.Now().UTC()
	dbLabel := labelData.Label{
		Created: now,
		Name:    labelName,
		Updated: now,
	}

	err = t.repository.Create(err, &dbLabel)
	return t.getWithoutContextCheck(err, dbLabel.Id)
}

func (t *LabelService) getWithoutContextCheck(err error, id int) (labelModels.Label, error) {
	if err != nil {
		return labelModels.Label{}, err
	}

	dbLabel, err := t.repository.GetOne(err, id)
	return labelModels.Label{
		Id:   dbLabel.Id,
		Name: dbLabel.Name,
	}, err
}

func (t *LabelService) modifyInternal(err error, currentManager labelModels.ManagerContext, lebelName string) (labelModels.Label, error) {
	if err != nil {
		return labelModels.Label{}, err
	}

	dbLabel, err := t.repository.GetOne(err, currentManager.LabelId)

	dbLabel.Name = lebelName
	err = t.repository.Update(err, &dbLabel)
	return t.getWithoutContextCheck(err, dbLabel.Id)
}
