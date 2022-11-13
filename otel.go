package main

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func buildProvider(ctx context.Context) (*tracesdk.TracerProvider, *metric.MeterProvider, error) {
	// metrics
	m_exp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	mp := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(m_exp)))
	defer func() {
		if err = mp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	global.SetMeterProvider(mp)

	// tracer
	t_exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		panic(err)
	}
	tp := tracesdk.NewTracerProvider(tracesdk.WithBatcher(t_exp))
	defer func() {
		if err = tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	return tp, mp, nil
}
