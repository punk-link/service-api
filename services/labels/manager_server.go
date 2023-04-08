package labels

import labelModels "main/models/labels"

type ManagerServer interface {
	Add(currentManager labelModels.ManagerContext, manager labelModels.Manager) (labelModels.Manager, error)
	AddMaster(request labelModels.AddMasterManagerRequest) (labelModels.Manager, error)
	Get(currentManager labelModels.ManagerContext) ([]labelModels.Manager, error)
	GetContext(id int) (labelModels.ManagerContext, error)
	GetOne(currentManager labelModels.ManagerContext, id int) (labelModels.Manager, error)
	Modify(currentManager labelModels.ManagerContext, manager labelModels.Manager, id int) (labelModels.Manager, error)
}
