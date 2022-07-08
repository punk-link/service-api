package organizations

import (
	"main/data"
	organizationData "main/data/organizations"
	"main/models/organizations"
	requests "main/requests/organizations"
	"main/services/helpers"
	"time"

	"github.com/rs/zerolog/log"
)

func AddManager(currentManager organizations.Manager, manager organizations.Manager) (organizations.Manager, error) {
	trimmedName, err := validateAndTrimName(manager.Name)
	if err != nil {
		return organizations.Manager{}, err
	}

	now := time.Now().UTC()
	dbManager := organizationData.Manager{
		Created:        now,
		Name:           trimmedName,
		OrganizationId: currentManager.OrganizationId,
		Updated:        now,
	}

	result := data.DB.Create(&dbManager)
	if result.Error != nil {
		log.Logger.Error().Err(result.Error).Msg(result.Error.Error())
		return organizations.Manager{}, result.Error
	}

	return GetManager(dbManager.Id)
}

func AddMasterManager(request requests.AddMasterManagerRequest) (organizations.Manager, error) {
	organization, err := AddOrganization(request.OrganizationName)
	if err != nil {
		return organizations.Manager{}, err
	}

	manager, err := AddManager(organizations.Manager{OrganizationId: organization.Id}, organizations.Manager{Name: request.Name})
	if err != nil {
		return organizations.Manager{}, err
	}

	return GetManager(manager.Id)
}

func GetManager(id int) (organizations.Manager, error) {
	dbManager, err := helpers.GetData[organizationData.Manager](id)
	if err != nil {
		return organizations.Manager{}, err
	}

	return toManager(dbManager), nil
}

func GetOrganizationManagers(organizationId int) []organizations.Manager {
	var dbManagers []organizationData.Manager
	data.DB.Where("organization_id = ?", organizationId).Find(&dbManagers)

	managers := make([]organizations.Manager, len(dbManagers))
	for i, manager := range dbManagers {
		managers[i] = toManager(manager)
	}

	return managers
}

func ModifyManager(manager organizations.Manager, id int) (organizations.Manager, error) {
	if err := validateId(manager.Id, id); err != nil {
		return organizations.Manager{}, err
	}

	trimmedName, err := validateAndTrimName(manager.Name)
	if err != nil {
		return organizations.Manager{}, err
	}

	dbManager, err := helpers.GetData[organizationData.Manager](id)
	if err != nil {
		return organizations.Manager{}, err
	}

	dbManager.Name = trimmedName
	dbManager.Updated = time.Now().UTC()
	data.DB.Save(&dbManager)

	return GetManager(id)
}

func toManager(source organizationData.Manager) organizations.Manager {
	return organizations.Manager{
		Id:             source.Id,
		Name:           source.Name,
		OrganizationId: source.OrganizationId,
	}
}
