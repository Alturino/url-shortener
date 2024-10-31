package log

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
)

func newHttpExporter(c context.Context) (trace.SpanExporter, error) {
	return otlptracehttp.New(c)
}
