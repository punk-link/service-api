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
	ctx.JSON(http.StatusBadRequest, gin.H{
		"message": reason,
	})
}

func Ok[T any](ctx *gin.Context, result T) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func UnprocessableEntity(ctx *gin.Context, binder error) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": binder.Error(),
	})
}
