package labels

import (
	"main/data"
	labelData "main/data/labels"
	"main/helpers"
	"main/models/labels"
	"main/services/common"
	validator "main/services/labels/validators"
	"main/services/platforms/spotify"
	"strings"
	"time"

	"github.com/samber/do"
)

type LabelService struct {
	logger         *common.Logger
	spotifyService *spotify.SpotifyService
}

func ConstructLabelService(injector *do.Injector) (*LabelService, error) {
	logger := do.MustInvoke[*common.Logger](injector)
	spotifyService := do.MustInvoke[*spotify.SpotifyService](injector)

	return &LabelService{
		logger:         logger,
		spotifyService: spotifyService,
	}, nil
}

func (t *LabelService) AddLabel(labelName string) (labels.Label, error) {
	trimmedName := strings.TrimSpace(labelName)
	err := validator.NameNotEmpty(trimmedName)

	return t.addLabelInternal(err, trimmedName)
}

func (t *LabelService) GetLabel(currentManager labels.ManagerContext, id int) (labels.Label, error) {
	err := validator.CurrentManagerBelongsToLabel(currentManager, id)
	return t.getLabelWithoutContextCheck(err, id)
}

func (t *LabelService) ModifyLabel(currentManager labels.ManagerContext, label labels.Label, id int) (labels.Label, error) {
	err1 := validator.CurrentManagerBelongsToLabel(currentManager, id)
	err2 := validator.IdConsistsOverRequest(label.Id, id)

	trimmedName := strings.TrimSpace(label.Name)
	err3 := validator.NameNotEmpty(trimmedName)

	return t.modifyLabelInternal(helpers.AccumulateErrors(err1, err2, err3), currentManager, trimmedName)
}

func (t *LabelService) addLabelInternal(err error, labelName string) (labels.Label, error) {
	if err != nil {
		return labels.Label{}, err
	}

	now := time.Now().UTC()
	dbLabel := labelData.Label{
		Created: now,
		Name:    labelName,
		Updated: now,
	}

	err = data.DB.Create(&dbLabel).Error
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return t.getLabelWithoutContextCheck(err, dbLabel.Id)
}

func (t *LabelService) getLabelWithoutContextCheck(err error, id int) (labels.Label, error) {
	if err != nil {
		return labels.Label{}, err
	}

	var dbLabel labelData.Label
	err = data.DB.First(&dbLabel, id).Error
	if err != nil {
		return labels.Label{}, err
	}

	return labels.Label{
		Id:   dbLabel.Id,
		Name: dbLabel.Name,
	}, nil
}

func (t *LabelService) modifyLabelInternal(err error, currentManager labels.ManagerContext, lebelName string) (labels.Label, error) {
	if err != nil {
		return labels.Label{}, err
	}

	var dbLabel labelData.Label
	err = data.DB.First(&dbLabel, currentManager.LabelId).Error
	if err != nil {
		return labels.Label{}, err
	}

	dbLabel.Name = lebelName
	dbLabel.Updated = time.Now().UTC()
	err = data.DB.Save(&dbLabel).Error

	return t.getLabelWithoutContextCheck(err, dbLabel.Id)
}
