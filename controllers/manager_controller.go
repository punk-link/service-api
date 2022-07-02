package controllers

import (
	"main/models/organizations"
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

	result, err := service.AddManager(manager)
	OkOrBadRequest(ctx, result, err)
}

func GetManager(ctx *gin.Context) {
	managerId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	result, err := service.GetManager(managerId)
	OkOrBadRequest(ctx, result, err)
}

func ModifyManager(ctx *gin.Context) {
	managerId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	var manager organizations.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	result, err := service.ModifyManager(manager, managerId)
	OkOrBadRequest(ctx, result, err)
}
