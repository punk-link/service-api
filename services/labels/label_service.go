package labels

import (
	"main/data"
	labelData "main/data/labels"
	"main/helpers"
	"main/models/labels"
	"main/services/common"
	validator "main/services/labels/validators"
	"main/services/spotify"
	"strings"
	"time"
)

type LabelService struct {
	logger         *common.Logger
	spotifyService *spotify.SpotifyService
}

func BuildLabelService(logger *common.Logger, spotifyService *spotify.SpotifyService) *LabelService {
	return &LabelService{
		logger:         logger,
		spotifyService: spotifyService,
	}
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
		service.logger.LogError(result.Error, result.Error.Error())
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
