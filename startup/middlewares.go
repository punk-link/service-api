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
	"/health":  true,
	"/metrics": true,
}
