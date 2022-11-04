package api

import (
	base "main/controllers"
	"main/models/labels"
	service "main/services/labels"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type LabelController struct {
	labelService   *service.LabelService
	managerService *service.ManagerService
}

func NewLabelController(injector *do.Injector) (*LabelController, error) {
	labelService := do.MustInvoke[*service.LabelService](injector)
	managerService := do.MustInvoke[*service.ManagerService](injector)

	return &LabelController{
		labelService:   labelService,
		managerService: managerService,
	}, nil
}

func (t *LabelController) Get(ctx *gin.Context) {
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

	result, err := t.labelService.GetLabel(currentManager, id)
	base.OkOrBadRequest(ctx, result, err)
}

func (t *LabelController) Modify(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	var label labels.Label
	if err := ctx.ShouldBindJSON(&label); err != nil {
		base.UnprocessableEntity(ctx, err)
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		base.NotFound(ctx, err.Error())
		return
	}

	result, err := t.labelService.ModifyLabel(currentManager, label, id)
	base.OkOrBadRequest(ctx, result, err)
}
