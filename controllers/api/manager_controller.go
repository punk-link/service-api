package api

import (
	base "main/controllers"
	"main/models/labels"
	service "main/services/labels"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ManagerController struct {
	managerService *service.ManagerService
}

func ConstructManagerController(managerService *service.ManagerService) *ManagerController {
	return &ManagerController{
		managerService: managerService,
	}
}

func (controller *ManagerController) Add(ctx *gin.Context) {
	var manager labels.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		base.UnprocessableEntity(ctx, err)
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx)
	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := controller.managerService.Add(currentManager, manager)
	base.OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) AddMaster(ctx *gin.Context) {
	var request labels.AddMasterManagerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		base.UnprocessableEntity(ctx, err)
		return
	}

	result, err := controller.managerService.AddMaster(request)
	base.OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) GetOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx)
	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := controller.managerService.GetOne(currentManager, id)
	base.OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) Get(ctx *gin.Context) {
	currentManager, err := base.GetCurrentManagerContext(ctx)
	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := controller.managerService.Get(currentManager)
	base.OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) Modify(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	var manager labels.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		base.UnprocessableEntity(ctx, err)
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx)
	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := controller.managerService.Modify(currentManager, manager, id)
	base.OkOrBadRequest(ctx, result, err)
}
