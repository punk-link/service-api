package controllers

import (
	"main/models/managers"
	service "main/services/managers"

	"github.com/gin-gonic/gin"
)

func AddManager(ctx *gin.Context) {
	var manager managers.Manager
	if err := ctx.ShouldBindJSON(&manager); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	result, err := service.AddManager(manager)
	OkOrBadRequest(ctx, result, err)
}

func GetManager(ctx *gin.Context) {

}
