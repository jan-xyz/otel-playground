package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func buildProvider(ctx context.Context) (*tracesdk.TracerProvider, *metricsdk.MeterProvider, error) {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(service),
		attribute.String("environment", environment),
		attribute.Int64("ID", id),
	)

	// Create Metrics exporter
	m_exp, err := stdoutmetric.New()
	if err != nil {
		return nil, nil, fmt.Errorf("creating metric exporter: %w", err)
	}
	mp := metricsdk.NewMeterProvider(
		metricsdk.WithReader(metricsdk.NewPeriodicReader(m_exp)),
		// Record information about this application in a Resource.
		metricsdk.WithResource(res),
	)

	// Create Trace exporter
	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("creating tace exporter: %w", err)
	}

	return tp, mp, nil
}
