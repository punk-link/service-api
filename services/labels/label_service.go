package labels

import (
	"main/data"
	labelData "main/data/labels"
	"main/models/labels"
	"main/services/helpers"
	validator "main/services/labels/validators"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type LabelService struct{}

func NewLabelService() *LabelService {
	return &LabelService{}
}

func (service *LabelService) AddLabel(labelName string) (labels.Label, error) {
	trimmedName := strings.TrimSpace(labelName)
	if err := validator.NameNotEmpty(trimmedName); err != nil {
		return labels.Label{}, err
	}

	now := time.Now().UTC()
	dbLabel := labelData.Label{
		Created: now,
		Name:    trimmedName,
		Updated: now,
	}

	result := data.DB.Create(&dbLabel)
	if result.Error != nil {
		log.Logger.Error().Err(result.Error).Msg(result.Error.Error())
		return labels.Label{}, result.Error
	}

	return getLabelWithoutContextCheck(dbLabel.Id)
}

func (service *LabelService) GetLabel(currentManager labels.ManagerContext, id int) (labels.Label, error) {
	if err := validator.CurrentManagerBelongsToLabel(currentManager, id); err != nil {
		return labels.Label{}, err
	}

	return getLabelWithoutContextCheck(id)
}

func (service *LabelService) ModifyLabel(currentManager labels.ManagerContext, label labels.Label, id int) (labels.Label, error) {
	if err := validator.CurrentManagerBelongsToLabel(currentManager, id); err != nil {
		return labels.Label{}, err
	}

	if err := validator.IdConsistsOverRequest(label.Id, id); err != nil {
		return labels.Label{}, err
	}

	trimmedName := strings.TrimSpace(label.Name)
	if err := validator.NameNotEmpty(trimmedName); err != nil {
		return labels.Label{}, err
	}

	dbLabel, err := helpers.GetEntity[labelData.Label](label.Id)
	if err != nil {
		return labels.Label{}, err
	}

	dbLabel.Name = trimmedName
	dbLabel.Updated = time.Now().UTC()
	data.DB.Save(&dbLabel)

	return service.GetLabel(currentManager, label.Id)
}

func getLabelWithoutContextCheck(id int) (labels.Label, error) {
	dbLabel, err := helpers.GetEntity[labelData.Label](id)
	if err != nil {
		return labels.Label{}, err
	}

	return labels.Label{
		Id:   dbLabel.Id,
		Name: dbLabel.Name,
	}, nil
}
