package validators

import (
	"errors"
	labelModels "main/models/labels"
)

func CurrentManagerBelongsToLabel(currentManager labelModels.ManagerContext, targetLabelId int) error {
	if currentManager.LabelId != targetLabelId {
		return errors.New("the current mamber isn't belong to the target label")
	}

	return nil
}

func IdConsistsOverRequest(entityId int, queryId int) error {
	if entityId != queryId {
		return errors.New("inconsistent ids in a body and a query")
	}

	return nil
}

func NameNotEmpty(name string) error {
	if len(name) == 0 {
		return errors.New("name must be provided")
	}

	return nil
}
