package controllers

import (
	"main/models/labels"
	requests "main/requests/labels"
	service "main/services/labels"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ManagerController struct {
	managerService *service.ManagerService
}

func NewManagerController(managerService *service.ManagerService) *ManagerController {
	return &ManagerController{
		managerService: managerService,
	}
}

func (controller *ManagerController) AddManager(ctx *gin.Context) {
	var manager labels.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := controller.managerService.AddManager(currentManager, manager)
	OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) AddMasterManager(ctx *gin.Context) {
	var request requests.AddMasterManagerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	result, err := controller.managerService.AddMasterManager(request)
	OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) GetManager(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := controller.managerService.GetManager(currentManager, id)
	OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) GetManagers(ctx *gin.Context) {
	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result := controller.managerService.GetLabelManagers(currentManager)
	Ok(ctx, result)
}

func (controller *ManagerController) ModifyManager(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	var manager labels.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := controller.managerService.ModifyManager(currentManager, manager, id)
	OkOrBadRequest(ctx, result, err)
}
