package startup

import (
	"fmt"

	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _exceptionalRoutes[ctx.Request.URL.Path] {
			ctx.Next()
			return
		}

		ctx.Next()
		counterName := fmt.Sprintf("requests_total{path=%q, status_code=\"%d\"}", ctx.Request.URL.Path, ctx.Writer.Status())
		metrics.GetOrCreateCounter(counterName).Inc()
	}
}

var _exceptionalRoutes map[string]bool = map[string]bool{
	"/favicon.ico":                     true,
	"/health":                          true,
	"/metrics":                         true,
	"/metrics/api/v1/query":            true,
	"/metrics/api/v1/query_range":      true,
	"/metrics/api/v1/rules":            true,
	"/metrics/api/v1/series":           true,
	"/metrics/api/v1/status/buildinfo": true,
}
