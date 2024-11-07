package log

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type ShutdownFunc func(context.Context) error

func newPropagator() propagation.TextMapPropagator {
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	return propagator
}

func newTracerProvider(c context.Context) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracehttp.New(
		c,
		otlptracehttp.WithEndpoint("otel-collector:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		logger.Error().
			Err(err).
			Str(KeyProcess, "main").
			Msgf("failed creating traceExporter with error=%s", err.Error())
		return nil, err
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(5*time.Second)),
	)
	return traceProvider, nil
}

func newMeterProvider(c context.Context) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetrichttp.New(
		c,
		otlpmetrichttp.WithEndpoint("otel-collector:4318"),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		logger.Error().
			Err(err).
			Str(KeyProcess, "main").
			Msgf("failed creating metricExporter with error=%s", err.Error())
		return nil, err
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(metricExporter, metric.WithInterval(5*time.Second)),
		),
	)
	return meterProvider, nil
}

func newLoggerProvider() (*log.LoggerProvider, error) {
	logExporter, err := stdoutlog.New()
	if err != nil {
		logger.Error().
			Err(err).
			Str(KeyProcess, "main").
			Msgf("failed creating logExporter with error=%s", err.Error())
		return nil, err
	}
	loggerProvider := log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(logExporter)))
	if err != nil {
		logger.Fatal().
			Err(err).
			Str(KeyProcess, "main").
			Msgf("failed creating loggerProvider with error=%s", err.Error())
	}
	return loggerProvider, nil
}

func InitOtelSdk(c context.Context) (shutdown ShutdownFunc, err error) {
	shutdownFuncs := []ShutdownFunc{}

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(c))
	}

	logger.Info().Str(KeyProcess, "main").Msg("initializing otel propagator")
	propagator := newPropagator()
	otel.SetTextMapPropagator(propagator)
	logger.Info().Str(KeyProcess, "main").Msg("initialized otel propagator")

	logger.Info().Str(KeyProcess, "main").Msg("initializing otel propagator")
	traceProvider, err := newTracerProvider(c)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, traceProvider.Shutdown)
	otel.SetTracerProvider(traceProvider)
	logger.Info().Str(KeyProcess, "main").Msg("initialized otel traceProvider")

	logger.Info().Str(KeyProcess, "main").Msg("initializing otel meterProvider")
	meterProvider, err := newMeterProvider(c)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)
	logger.Info().Str(KeyProcess, "main").Msg("initialized otel meterProvider")

	logger.Info().Str(KeyProcess, "main").Msg("initializing otel loggerProvider")
	loggerProvider, err := newLoggerProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	logger.Info().Str(KeyProcess, "main").Msg("initialized otel loggerProvider")

	return
}
