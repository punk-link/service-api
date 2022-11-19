package startup

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
)

func metricsMiddleware(instrumentationName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _exceptionalRoutes[ctx.Request.URL.Path] {
			ctx.Next()
			return
		}

		ctx.Next()

		meter := global.MeterProvider().Meter(instrumentationName)
		requestCounter, _ := meter.SyncInt64().Counter("requests")

		attributes := []attribute.KeyValue{
			attribute.Key("path").String(ctx.Request.URL.Path),
			attribute.Key("status_code").Int(ctx.Writer.Status()),
		}
		requestCounter.Add(ctx, 1, attributes...)
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
