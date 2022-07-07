package organizations

type Manager struct {
	Id             int    `json:"id"`
	OrganizationId int    `json:"organizationId"`
	Name           string `json:"name"`
}
