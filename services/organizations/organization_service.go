package organizations

import (
	"main/data"
	organizationData "main/data/organizations"
	"main/models/organizations"
	"main/services/helpers"
	"time"

	"github.com/rs/zerolog/log"
)

func AddOrganization(organizationName string) (organizations.Organization, error) {
	trimmedName, err := validateAndTrimName(organizationName)
	if err != nil {
		return organizations.Organization{}, err
	}

	now := time.Now().UTC()
	dbOrganization := organizationData.Organization{
		Created: now,
		Name:    trimmedName,
		Updated: now,
	}

	result := data.DB.Create(&dbOrganization)
	if result.Error != nil {
		log.Logger.Error().Err(result.Error).Msg(result.Error.Error())
		return organizations.Organization{}, result.Error
	}

	return GetOrganization(dbOrganization.Id)
}

func GetOrganization(id int) (organizations.Organization, error) {
	dbOrganization, err := helpers.GetData[organizationData.Organization](id)
	if err != nil {
		return organizations.Organization{}, err
	}

	return organizations.Organization{
		Id:   dbOrganization.Id,
		Name: dbOrganization.Name,
	}, nil
}

func ModifyOrganization(organization organizations.Organization, id int) (organizations.Organization, error) {
	if err := validateId(organization.Id, id); err != nil {
		return organizations.Organization{}, err
	}

	trimmedName, err := validateAndTrimName(organization.Name)
	if err != nil {
		return organizations.Organization{}, err
	}

	dbOrganization, err := helpers.GetData[organizationData.Organization](organization.Id) //getOrganizationData(id)
	if err != nil {
		return organizations.Organization{}, err
	}

	dbOrganization.Name = trimmedName
	dbOrganization.Updated = time.Now().UTC()
	data.DB.Save(&dbOrganization)

	return GetOrganization(organization.Id)
}
