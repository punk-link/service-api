package controllers

import (
	"main/models/organizations"
	service "main/services/organizations"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOrganization(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	result, err := service.GetOrganization(id)
	OkOrBadRequest(ctx, result, err)
}

func ModifyOrganization(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	var organization organizations.Organization
	if err := ctx.ShouldBindJSON(&organization); err != nil {
		UnprocessableEntity(ctx, err)
		return
	}

	result, err := service.ModifyOrganization(organization, id)
	OkOrBadRequest(ctx, result, err)
}
