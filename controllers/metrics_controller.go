package controllers

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type MetricsController struct {
}

func NewMetricsController(injector *do.Injector) (*MetricsController, error) {
	return &MetricsController{}, nil
}

func (t *MetricsController) GetMetrics(ctx *gin.Context) {
	metrics.WritePrometheus(ctx.Writer, true)
}
