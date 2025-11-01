package monitoring

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpMetrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type MetricProvider struct {
	httpMiddleware middleware.Middleware
	stats.ConnStats

	metric.MeterProvider
}

func NewMetricProvider(name Name, version Version) (*MetricProvider, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	provider := sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(exporter),
		sdkMetric.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(string(name)),
			semconv.ServiceVersionKey.String(string(version)),
		)),
	)

	// TODO implement go.opentelemetry recorder
	httpMiddleware := middleware.New(middleware.Config{
		Recorder: httpMetrics.NewRecorder(httpMetrics.Config{}),
	})

	return &MetricProvider{
		MeterProvider:  provider,
		httpMiddleware: httpMiddleware,
	}, nil
}

func (m *MetricProvider) HTTPHandlerProvider(handleID string) func(http.Handler) http.Handler {
	return std.HandlerProvider(handleID, m.httpMiddleware)
}

func (m *MetricProvider) GRPCServerOption() grpc.ServerOption {
	return grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithMeterProvider(m)))
}

func (m *MetricProvider) Handler() http.Handler {
	return promhttp.Handler()
}
