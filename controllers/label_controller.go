package controllers

import (
	"main/models/labels"
	service "main/services/labels"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLabel(ctx *gin.Context) {
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

	result, err := service.GetLabel(id, currentManager)
	OkOrBadRequest(ctx, result, err)
}

func ModifyLabel(ctx *gin.Context) {
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

	result, err := service.ModifyLabel(label, id, currentManager)
	OkOrBadRequest(ctx, result, err)
}
