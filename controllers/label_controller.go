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

	result, err := service.GetLabel(id)
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

	result, err := service.ModifyLabel(label, id)
	OkOrBadRequest(ctx, result, err)
}
