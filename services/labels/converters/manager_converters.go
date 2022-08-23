package converters

import (
	data "main/data/labels"
	models "main/models/labels"
	"time"
)

func ToManager(source data.Manager) models.Manager {
	return models.Manager{
		Id:      source.Id,
		Name:    source.Name,
		LabelId: source.LabelId,
	}
}

func ToManagers(source []data.Manager) []models.Manager {
	managers := make([]models.Manager, len(source))
	for i, manager := range source {
		managers[i] = ToManager(manager)
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
