package tracing

import (
	"Web3-Telegram-Wallet-Bot/internal/config"
	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	traceI "go.opentelemetry.io/otel/trace"
)

func NewTracerProvider(ctx context.Context, config *config.Tracing) (*trace.TracerProvider, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	exporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(config.Endpoint),
		otlptracehttp.WithHeaders(headers),
		otlptracehttp.WithInsecure(),
	),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create exporter")
	}
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("web3-wallet"),
			)))
	otel.SetTracerProvider(provider)

	return provider, nil
}

func NewMockTracer() (traceI.Tracer, *tracetest.SpanRecorder) {
	rec := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(rec),
	)
	tracer := tp.Tracer("test")
	return tracer, rec
}
