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

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := service.AddManager(manager, currentManager)
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

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := service.GetManager(id, currentManager)
	OkOrBadRequest(ctx, result, err)
}

func GetManagers(ctx *gin.Context) {
	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result := service.GetLabelManagers(currentManager)
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

	currentManager, err := getCurrentManagerContext(ctx)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	result, err := service.ModifyManager(manager, id, currentManager)
	OkOrBadRequest(ctx, result, err)
}
