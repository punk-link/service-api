package api

import (
	base "main/controllers"
	"main/models/labels"
	service "main/services/labels"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type ManagerController struct {
	managerService *service.ManagerService
}

func ConstructManagerController(injector *do.Injector) (*ManagerController, error) {
	managerService := do.MustInvoke[*service.ManagerService](injector)

	return &ManagerController{
		managerService: managerService,
	}, nil
}

func (t *ManagerController) Add(ctx *gin.Context) {
	var manager labels.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		base.UnprocessableEntity(ctx, err)
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := t.managerService.Add(currentManager, manager)
	base.OkOrBadRequest(ctx, result, err)
}

func (t *ManagerController) AddMaster(ctx *gin.Context) {
	var request labels.AddMasterManagerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		base.UnprocessableEntity(ctx, err)
		return
	}

	result, err := t.managerService.AddMaster(request)
	base.OkOrBadRequest(ctx, result, err)
}

func (t *ManagerController) GetOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := t.managerService.GetOne(currentManager, id)
	base.OkOrBadRequest(ctx, result, err)
}

func (t *ManagerController) Get(ctx *gin.Context) {
	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := t.managerService.Get(currentManager)
	base.OkOrBadRequest(ctx, result, err)
}

func (t *ManagerController) Modify(ctx *gin.Context) {
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

	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := t.managerService.Modify(currentManager, manager, id)
	base.OkOrBadRequest(ctx, result, err)
}
