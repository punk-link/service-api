package controllers

import (
	"main/models/labels"
	"main/services/common"
	service "main/services/labels"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCurrentManagerContext(ctx *gin.Context) (labels.ManagerContext, error) {
	headerValue := ctx.Request.Header["X-User-Id"][0]
	managerId, err := strconv.Atoi(headerValue)
	if err != nil {
		return labels.ManagerContext{}, err
	}

	// TODO: add an injection maybe
	managerService := service.BuildManagerService(&service.LabelService{}, &common.Logger{})
	return managerService.GetManagerContext(managerId)
}
