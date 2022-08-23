package controllers

import (
	"main/models/labels"
	service "main/services/labels"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LabelController struct {
	labelService *service.LabelService
}

func ConstructLabelController(labelService *service.LabelService) *LabelController {
	return &LabelController{
		labelService: labelService,
	}
}

func (t *LabelController) GetLabel(ctx *gin.Context) {
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

	result, err := t.labelService.GetLabel(currentManager, id)
	OkOrBadRequest(ctx, result, err)
}

func (t *LabelController) ModifyLabel(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	var label labels.Label
	if err := ctx.ShouldBindJSON(&label); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := t.labelService.ModifyLabel(currentManager, label, id)
	OkOrBadRequest(ctx, result, err)
}
