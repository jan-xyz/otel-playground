package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func setupMetrics(ctx context.Context, res *resource.Resource) *metric.MeterProvider {
	// communicate on localhost to the ADOT collector
	m_exp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint("0.0.0.0:4317"),
	)
	if err != nil {
		panic(err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(m_exp)),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	return mp
}
