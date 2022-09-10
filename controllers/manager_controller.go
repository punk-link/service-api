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

func ConstructManagerController(managerService *service.ManagerService) *ManagerController {
	return &ManagerController{
		managerService: managerService,
	}
}

func (controller *ManagerController) Add(ctx *gin.Context) {
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

	result, err := controller.managerService.Add(currentManager, manager)
	OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) AddMaster(ctx *gin.Context) {
	var request requests.AddMasterManagerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	result, err := controller.managerService.AddMaster(request)
	OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) GetOne(ctx *gin.Context) {
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

	result, err := controller.managerService.GetOne(currentManager, id)
	OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) Get(ctx *gin.Context) {
	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := controller.managerService.Get(currentManager)
	OkOrBadRequest(ctx, result, err)
}

func (controller *ManagerController) Modify(ctx *gin.Context) {
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

	result, err := controller.managerService.Modify(currentManager, manager, id)
	OkOrBadRequest(ctx, result, err)
}
