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

func OkOrNotFoundTemplate(ctx *gin.Context, templateName string, data map[string]any, err error) {
	if err != nil {
		NotFoundTemplate(ctx, err.Error())
		return
	}

	OkTemplate(ctx, templateName, data)
}

func NotFound(ctx *gin.Context, reason string) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"message": reason,
	})
}

func NotFoundTemplate(ctx *gin.Context, reason string) {
	ctx.HTML(http.StatusNotFound, "global/404.tmpl", gin.H{
		"Error":     reason,
		"PageTitle": "404",
	})
}

func UnprocessableEntity(ctx *gin.Context, binder error) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": binder.Error(),
	})
}
