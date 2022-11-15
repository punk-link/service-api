package startup

import (
	startupModels "main/models/startup"

	consulClient "github.com/punk-link/consul-client"
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

func configureOpenTelemetry(logger logger.Logger, consul *consulClient.ConsulClient, options *startupModels.StartupOptions) {
	configureTracing(logger, consul, options)
	configureMetrics(logger)
}

func configureMetrics(logger logger.Logger) {
	exporter, err := prometheus.New()
	logOpenTelemetryExceptionIfAny(logger, err)

	metricProvider := metricSdk.NewMeterProvider(metricSdk.WithReader(exporter))
	global.SetMeterProvider(metricProvider)
}

func configureTracing(logger logger.Logger, consul *consulClient.ConsulClient, options *startupModels.StartupOptions) {
	jaegerSettingsValues, err := consul.Get("JaegerSettings")
	if err != nil {
		logger.LogInfo("Jaeger settings is empty")
		return
	}

	jaegerSettings := jaegerSettingsValues.(map[string]any)
	endpoint := jaegerSettings["Endpoint"].(string)
	if endpoint == "" {
		logger.LogInfo("Jaeger endpoint is empty")
		return
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
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
