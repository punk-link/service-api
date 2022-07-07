package organizations

import (
	"errors"
	"strings"
)

func validateId(entityId int, queryId int) error {
	if entityId != queryId {
		return errors.New("inconsistent ids in a body and a query")
	}

	return nil
}

func validateAndTrimName(name string) (string, error) {
	trimmedName := strings.TrimSpace(name)
	if len(trimmedName) == 0 {
		return "", errors.New("name must be provided")
	}

	return trimmedName, nil
}
