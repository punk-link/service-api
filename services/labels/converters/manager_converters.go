package converters

import (
	data "main/data/labels"
	models "main/models/labels"
	"time"
)

func ToManager(err error, source data.Manager) (models.Manager, error) {
	if err != nil {
		return models.Manager{}, err
	}

	return models.Manager{
		Id:      source.Id,
		Name:    source.Name,
		LabelId: source.LabelId,
	}, nil
}

func ToManagers(source []data.Manager) []models.Manager {
	managers := make([]models.Manager, len(source))
	for i, manager := range source {
		managers[i], _ = ToManager(nil, manager)
	}

	return managers
}

func ToDbManager(managerName string, labelId int) data.Manager {
	now := time.Now().UTC()

	return data.Manager{
		Created: now,
		Name:    managerName,
		LabelId: labelId,
		Updated: now,
	}
}
