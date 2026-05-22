package observability

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// TracerProvider wraps the OpenTelemetry TracerProvider with
// ManPaSik-specific defaults (Jaeger OTLP endpoint, service resource).
type TracerProvider struct {
	provider *sdktrace.TracerProvider
}

// InitTracer creates and registers a global OTLP-based TracerProvider.
// The OTLP HTTP exporter targets the Jaeger collector (or any OTLP-compatible
// endpoint) configured via OTEL_EXPORTER_OTLP_ENDPOINT env var, falling back
// to http://jaeger:4318.
func InitTracer(serviceName string) (*TracerProvider, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "jaeger:4318"
	}

	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("otlp exporter 생성 실패: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(getVersion()),
			attribute.String("environment", getEnvironment()),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("resource 생성 실패: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(getSamplingRate()))),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &TracerProvider{provider: tp}, nil
}

// Shutdown gracefully flushes and shuts down the tracer provider.
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	if tp == nil || tp.provider == nil {
		return nil
	}
	return tp.provider.Shutdown(ctx)
}

// Tracer returns a named tracer from the underlying provider.
func (tp *TracerProvider) Tracer(name string) trace.Tracer {
	return tp.provider.Tracer(name)
}

// getVersion returns the service version from environment variable.
func getVersion() string {
	if v := os.Getenv("SERVICE_VERSION"); v != "" {
		return v
	}
	return "dev"
}

// getEnvironment returns the deployment environment.
func getEnvironment() string {
	if e := os.Getenv("ENVIRONMENT"); e != "" {
		return e
	}
	return "development"
}

// getSamplingRate returns the trace sampling rate (0.0 ~ 1.0).
func getSamplingRate() float64 {
	env := getEnvironment()
	switch env {
	case "production":
		return 0.1
	case "staging":
		return 0.5
	default:
		return 1.0
	}
}
