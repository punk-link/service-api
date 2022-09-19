package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(ctx *gin.Context, reason string) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"message": reason,
	})
}

func InternalServerError(ctx *gin.Context, reason string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"message": reason,
	})
}

func Ok[T any](ctx *gin.Context, result T) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func OkTemplate(ctx *gin.Context, templateName string, data map[string]any) {
	ctx.HTML(http.StatusOK, templateName, data)
}

func OkOrBadRequest[T any](ctx *gin.Context, result T, err error) {
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	Ok(ctx, result)
}

func NotFound(ctx *gin.Context, reason string) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"message": reason,
	})
}

func UnprocessableEntity(ctx *gin.Context, binder error) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": binder.Error(),
	})
}
