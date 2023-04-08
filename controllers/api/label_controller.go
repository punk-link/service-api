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

type LabelController struct {
	labelService   service.LabelServer
	managerService service.ManagerServer
}

func NewLabelController(injector *do.Injector) (*LabelController, error) {
	labelService := do.MustInvoke[service.LabelServer](injector)
	managerService := do.MustInvoke[service.ManagerServer](injector)

	return &LabelController{
		labelService:   labelService,
		managerService: managerService,
	}, nil
}

func (t *LabelController) Get(ctx *gin.Context) {
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

	result, err := t.labelService.GetOne(currentManager, id)
	templates.OkOrBadRequest(ctx, result, err)
}

func (t *LabelController) Modify(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		templates.BadRequest(ctx, err.Error())
		return
	}

	var label labels.Label
	if err := ctx.ShouldBindJSON(&label); err != nil {
		templates.UnprocessableEntity(ctx, err)
		return
	}

	currentManager, err := base.GetCurrentManagerContext(ctx, t.managerService)
	if err != nil {
		templates.NotFound(ctx, err.Error())
		return
	}

	result, err := t.labelService.Modify(currentManager, label, id)
	templates.OkOrBadRequest(ctx, result, err)
}
