package controllers

import (
	"main/models/labels"
	services "main/services/labels"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCurrentManagerContext(ctx *gin.Context, service services.ManagerServer) (labels.ManagerContext, error) {
	headerValue := ctx.Request.Header["X-User-Id"][0]
	managerId, err := strconv.Atoi(headerValue)
	if err != nil {
		return labels.ManagerContext{}, err
	}

	return service.GetContext(managerId)
}
