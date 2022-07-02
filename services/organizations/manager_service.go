package organizations

import (
	"errors"
	"main/data"
	dataOrganizations "main/data/organizations"
	"main/models/organizations"
	"strings"

	"github.com/rs/zerolog/log"
)

func AddManager(manager organizations.Manager) (organizations.Manager, error) {
	if len(manager.Name) == 0 {
		return organizations.Manager{}, errors.New("manager's name must be provided")
	}

	dbManager := dataOrganizations.Manager{
		Name: strings.TrimSpace(manager.Name),
	}

	result := data.DB.Create(&dbManager)
	if result.Error != nil {
		log.Logger.Error().Err(result.Error).Msg(result.Error.Error())
		return organizations.Manager{}, result.Error
	}

	return organizations.Manager{
		Id:   dbManager.Id,
		Name: dbManager.Name,
	}, nil
}

func GetManager(managerId int) (organizations.Manager, error) {
	var dbManager dataOrganizations.Manager
	result := data.DB.First(&dbManager, managerId)

	if result.RowsAffected != 1 {
		if result.Error != nil {
			return organizations.Manager{}, result.Error
		}

		return organizations.Manager{}, errors.New("no managers found")
	}

	return organizations.Manager{
		Id:   dbManager.Id,
		Name: dbManager.Name,
	}, nil
}

func ModifyManager(manager organizations.Manager, managerId int) (organizations.Manager, error) {
	var dbManager dataOrganizations.Manager
	result := data.DB.First(&dbManager, managerId)

	if result.RowsAffected != 1 {
		if result.Error != nil {
			return organizations.Manager{}, result.Error
		}

		return organizations.Manager{}, errors.New("no managers found")
	}

	dbManager.Name = manager.Name
	data.DB.Save(&dbManager)

	return organizations.Manager{
		Id:   dbManager.Id,
		Name: dbManager.Name,
	}, nil
}
