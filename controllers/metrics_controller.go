package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/samber/do"
)

type MetricsController struct {
}

func NewMetricsController(injector *do.Injector) (*MetricsController, error) {
	return &MetricsController{}, nil
}

func (t *MetricsController) GetMetrics(ctx *gin.Context) {
	handler := promhttp.Handler()
	handler.ServeHTTP(ctx.Writer, ctx.Request)
}
