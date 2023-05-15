package repositories

import labelData "main/data/labels"

type LabelRepository interface {
	Create(err error, dbLabel *labelData.Label) error
	GetOne(err error, id int) (labelData.Label, error)
	Update(err error, dbLabel *labelData.Label) error
}
