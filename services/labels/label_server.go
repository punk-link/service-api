package labels

import labelModels "main/models/labels"

type LabelServer interface {
	Add(labelName string) (labelModels.Label, error)
	GetOne(currentManager labelModels.ManagerContext, id int) (labelModels.Label, error)
	Modify(currentManager labelModels.ManagerContext, label labelModels.Label, id int) (labelModels.Label, error)
}
