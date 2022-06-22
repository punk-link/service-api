package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckStatus(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
