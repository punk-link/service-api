package controllers

import (
	"main/models/labels"
	requests "main/requests/labels"
	service "main/services/labels"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddManager(ctx *gin.Context) {
	var manager labels.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	// TODO: add the current manager
	result, err := service.AddManager(labels.Manager{
		Id:      1,
		LabelId: 1,
	}, manager)
	OkOrBadRequest(ctx, result, err)
}

func AddMasterManager(ctx *gin.Context) {
	var request requests.AddMasterManagerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	result, err := service.AddMasterManager(request)
	OkOrBadRequest(ctx, result, err)
}

func GetManager(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	result, err := service.GetManager(id)
	OkOrBadRequest(ctx, result, err)
}

func GetManagers(ctx *gin.Context) {
	// TODO: put an actual manager here
	result := service.GetLabelManagers(1)
	Ok(ctx, result)
}

func ModifyManager(ctx *gin.Context) {
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

	result, err := service.ModifyManager(manager, id)
	OkOrBadRequest(ctx, result, err)
}
