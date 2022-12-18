package converters

import (
	data "main/data/labels"
	labelModels "main/models/labels"
	"time"
)

func ToManager(err error, source data.Manager) (labelModels.Manager, error) {
	if err != nil {
		return labelModels.Manager{}, err
	}

	return labelModels.Manager{
		Id:      source.Id,
		Name:    source.Name,
		LabelId: source.LabelId,
	}, nil
}

func ToManagers(source []data.Manager) []labelModels.Manager {
	managers := make([]labelModels.Manager, len(source))
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
