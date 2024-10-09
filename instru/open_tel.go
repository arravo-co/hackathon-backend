package instru

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

var (
	tracerProvider *trace.TracerProvider
)

type SetupOtel struct {
	*trace.TracerProvider
	*log.LoggerProvider
	*metric.MeterProvider
}

func Setup(ctx context.Context, opt *SetupOtel) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error
	shutdown = func(c context.Context) error {
		var err error
		for _, v := range shutdownFuncs {
			err = errors.Join(err, v(c))
		}
		shutdownFuncs = nil
		return err
	}
	/*
		handleErr := func(inErr error) {
			err = errors.Join(inErr, shutdown(ctx))
		}
	*/
	/*
		tracerProvider, err := NewTraceProvider()
		if err != nil {
			handleErr(err)
			return
		}

		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)*/

	shutdownFuncs = append(shutdownFuncs, opt.TracerProvider.Shutdown)
	//otel.SetTracerProvider(tracerProvider)
	//otel.SetTracerProvider(opt.TracerProvider)

	/*
		meterProvider, err := NewMeterProvider()
		if err != nil {
			handleErr(err)
			return
		}

		shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)*/
	shutdownFuncs = append(shutdownFuncs, opt.MeterProvider.Shutdown)
	otel.SetMeterProvider(opt.MeterProvider)

	/*
		loggerProvider, err := NewLoggerProvider()
		if err != nil {
			handleErr(err)
			return
		}

		shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
		global.SetLoggerProvider(loggerProvider)
	*/
	shutdownFuncs = append(shutdownFuncs, opt.LoggerProvider.Shutdown)
	global.SetLoggerProvider(opt.LoggerProvider)
	return
}

func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func NewTraceProvider() (*trace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	/*
		otlpExporter, err := NewOTLPTraceExporter(context.Background())
		if err != nil {
			return nil, err
		}
	*/
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Minute*30)),
		//trace.WithBatcher(otlpExporter, trace.WithBatchTimeout(time.Second)),
	)
	return traceProvider, nil
}

func NewMeterProvider() (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	/*	otlpMetricExporter, err := NewOTLPMetricsExporter(context.Background())
		if err != nil {
			return nil, err
		}
	*/
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			metric.WithInterval(30*time.Minute))),
		/*
			metric.WithReader(metric.NewPeriodicReader(otlpMetricExporter,
				metric.WithInterval(3*time.Second)),
				),*/
	)
	return meterProvider, nil
}

func NewLoggerProvider() (*log.LoggerProvider, error) {
	logExporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}
	/*
		otlpExporter, err := NewOTLPLogExporter(context.Background())
		if err != nil {
			return nil, err
		}
	*/
	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	//	log.WithProcessor(log.NewBatchProcessor(otlpExporter)),
	)
	return loggerProvider, nil
}

func NewOTLPTraceExporter(ctx context.Context) (trace.SpanExporter, error) {
	return otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("http://localhost:4318"))
}

func NewOTLPMetricsExporter(ctx context.Context) (metric.Exporter, error) {
	return otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpointURL("http://localhost:4318"))
}

func NewOTLPLogExporter(ctx context.Context) (log.Exporter, error) {
	return otlploghttp.New(ctx, otlploghttp.WithEndpoint("http://localhost:4318"))
}
