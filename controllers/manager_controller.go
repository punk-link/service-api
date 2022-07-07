package controllers

import (
	"main/models/organizations"
	requests "main/requests/organizations"
	service "main/services/organizations"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddManager(ctx *gin.Context) {
	var manager organizations.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	// TODO: add the current manager
	result, err := service.AddManager(organizations.Manager{}, manager)
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

func ModifyManager(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	var manager organizations.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	result, err := service.ModifyManager(manager, id)
	OkOrBadRequest(ctx, result, err)
}
