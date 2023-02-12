package main

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

func setupTracing(ctx context.Context, res *resource.Resource) *trace.TracerProvider {
	// communicate on localhost to the ADOT collector
	// required to generate xray compatible IDs
	t_exp, err := otlptracegrpc.New(ctx,

		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("0.0.0.0:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		panic(err)
	}
	tp := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(t_exp),
		trace.WithSampler(trace.AlwaysSample()),

		trace.WithIDGenerator(xray.NewIDGenerator()),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(xray.Propagator{})
	return tp
}
