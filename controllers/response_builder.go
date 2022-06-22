package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OkOrBadRequest[T any](ctx *gin.Context, result T, err error) {
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	Ok(ctx, result)
}

func BadRequest(ctx *gin.Context, reason string) {
	status := http.StatusBadRequest

	ctx.JSON(status, gin.H{
		"message": reason,
		"status":  status,
	})
}

func Ok[T any](ctx *gin.Context, result T) {
	status := http.StatusOK

	ctx.JSON(status, gin.H{
		"data":   result,
		"status": status,
	})
}

func UnprocessableEntity(ctx *gin.Context, binder error) {
	status := http.StatusUnprocessableEntity

	ctx.JSON(status, gin.H{
		"message": binder.Error(),
		"status":  status,
	})
}
