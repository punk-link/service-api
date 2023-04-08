package api

import (
	"main/models/labels"
	service "main/services/labels"
	"strconv"

	base "main/controllers"

	templates "github.com/punk-link/gin-generic-http-templates"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type ManagerController struct {
	managerService service.ManagerServer
}

func NewManagerController(injector *do.Injector) (*ManagerController, error) {
	managerService := do.MustInvoke[service.ManagerServer](injector)

	return &ManagerController{
		managerService: managerService,
	}, nil
}

func (t *ManagerController) Add(ctx *gin.Context) {
	var manager labels.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		templates.UnprocessableEntity(ctx, err)
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		templates.NotFound(ctx, err.Error())
		return
	}

	result, err := t.managerService.Add(currentManager, manager)
	templates.OkOrBadRequest(ctx, result, err)
}

func (t *ManagerController) AddMaster(ctx *gin.Context) {
	var request labels.AddMasterManagerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		templates.UnprocessableEntity(ctx, err)
		return
	}

	result, err := t.managerService.AddMaster(request)
	templates.OkOrBadRequest(ctx, result, err)
}

func (t *ManagerController) GetOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		templates.BadRequest(ctx, err.Error())
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		templates.NotFound(ctx, err.Error())
		return
	}

	result, err := t.managerService.GetOne(currentManager, id)
	templates.OkOrBadRequest(ctx, result, err)
}

func (t *ManagerController) Get(ctx *gin.Context) {
	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		templates.NotFound(ctx, err.Error())
		return
	}

	result, err := t.managerService.Get(currentManager)
	templates.OkOrBadRequest(ctx, result, err)
}

func (t *ManagerController) Modify(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		templates.BadRequest(ctx, err.Error())
		return
	}

	var manager labels.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		templates.UnprocessableEntity(ctx, err)
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		templates.NotFound(ctx, err.Error())
		return
	}

	result, err := t.managerService.Modify(currentManager, manager, id)
	templates.OkOrBadRequest(ctx, result, err)
}
