package startup

import (
	startupModels "main/models/startup"

	"github.com/punk-link/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	metricSdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func configureOpenTelemetry(logger logger.Logger, options *startupModels.StartupOptions) {
	configureTracing(logger, options)
	configureMetrics(logger)
}

func configureMetrics(logger logger.Logger) {
	exporter, err := prometheus.New()
	logOpenTelemetryExceptionIfAny(logger, err)

	metricProvider := metricSdk.NewMeterProvider(metricSdk.WithReader(exporter))
	global.SetMeterProvider(metricProvider)
}

func configureTracing(logger logger.Logger, options *startupModels.StartupOptions) {
	if options.JaegerEndpoint == "" {
		logger.LogInfo("Jaeger endpoint is empty")
		return
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(options.JaegerEndpoint)))
	logOpenTelemetryExceptionIfAny(logger, err)

	traceProvider := traceSdk.NewTracerProvider(traceSdk.WithBatcher(exporter), traceSdk.WithResource(resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(options.ServiceName),
		attribute.String("environment", options.EnvironmentName),
	)))

	otel.SetTracerProvider(traceProvider)
}

func logOpenTelemetryExceptionIfAny(logger logger.Logger, err error) {
	if err == nil {
		return
	}

	logger.LogFatal(err, "OpenTelemetry error: %s", err)
}
