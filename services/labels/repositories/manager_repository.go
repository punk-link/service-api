package repositories

import labelData "main/data/labels"

type ManagerRepository interface {
	Create(err error, dbManager *labelData.Manager) error
	GetByLabelId(err error, labelId int) ([]labelData.Manager, error)
	GetOne(err error, id int) (labelData.Manager, error)
	Update(err error, dbManager *labelData.Manager) error
}
