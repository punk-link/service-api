package startup

import (
	dataStructures "main/data-structures"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
)

func metricsMiddleware(instrumentationName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _exceptionalRoutes.Contains(ctx.Request.URL.Path) {
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

var _exceptionalRoutes dataStructures.HashSet[string] = dataStructures.MakeHashSet([]string{
	"/favicon.ico",
	"/health",
	"/metrics",
	"/metrics/api/v1/query",
	"/metrics/api/v1/query_range",
	"/metrics/api/v1/rules",
	"/metrics/api/v1/series",
	"/metrics/api/v1/status/buildinfo",
})
