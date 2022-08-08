package labels

import (
	"main/data"
	labelData "main/data/labels"
	"main/models/labels"
	"main/services/helpers"
	"time"

	"github.com/rs/zerolog/log"
)

func AddLabel(labelName string) (labels.Label, error) {
	trimmedName, err := validateAndTrimName(labelName)
	if err != nil {
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

	return GetLabel(dbLabel.Id)
}

func GetLabel(id int) (labels.Label, error) {
	dbLabel, err := helpers.GetData[labelData.Label](id)
	if err != nil {
		return labels.Label{}, err
	}

	return labels.Label{
		Id:   dbLabel.Id,
		Name: dbLabel.Name,
	}, nil
}

func ModifyLabel(label labels.Label, id int) (labels.Label, error) {
	if err := validateId(label.Id, id); err != nil {
		return labels.Label{}, err
	}

	trimmedName, err := validateAndTrimName(label.Name)
	if err != nil {
		return labels.Label{}, err
	}

	dbLabel, err := helpers.GetData[labelData.Label](label.Id)
	if err != nil {
		return labels.Label{}, err
	}

	dbLabel.Name = trimmedName
	dbLabel.Updated = time.Now().UTC()
	data.DB.Save(&dbLabel)

	return GetLabel(label.Id)
}
