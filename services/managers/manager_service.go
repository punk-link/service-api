package managers

import (
	"errors"
	"main/data"
	dataManagers "main/data/managers"
	"main/models/managers"
	"strings"

	"github.com/rs/zerolog/log"
)

func AddManager(manager managers.Manager) (managers.Manager, error) {
	if len(manager.Name) == 0 {
		return managers.Manager{}, errors.New("manager's name must be provided")
	}

	dbManager := dataManagers.Manager{
		Name: strings.TrimSpace(manager.Name),
	}

	result := data.DB.Create(&dbManager)
	if result.Error != nil {
		log.Logger.Error().Err(result.Error).Msg(result.Error.Error())
		return managers.Manager{}, result.Error
	}

	return managers.Manager{
		Id:   dbManager.Id,
		Name: dbManager.Name,
	}, nil
}
