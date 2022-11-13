package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

var (
	service           = "otel-test"
	environment       = "dev"
	id          int64 = 12345
)

func main() {
	ctx := context.Background()
	buildLogger()

	res, err := resource.New(ctx, resource.WithAttributes(
		// sets the service correctly
		semconv.ServiceNameKey.String(service),
	))
	if err != nil {
		panic(err)
	}
	// metrics
	m_exp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure(), otlpmetricgrpc.WithEndpoint("0.0.0.0:4317"), otlpmetricgrpc.WithDialOption(grpc.WithBlock()))
	if err != nil {
		panic(err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(m_exp)),
		metric.WithResource(res),
	)
	defer func() {
		if err = mp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	global.SetMeterProvider(mp)

	// tracer
	t_exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("0.0.0.0:4317"), otlptracegrpc.WithDialOption(grpc.WithBlock()))
	if err != nil {
		panic(err)
	}
	// required to generate xray compatible IDs
	idg := xray.NewIDGenerator()
	tp := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(t_exp),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithIDGenerator(idg),
	)
	defer func() {
		if err = tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(xray.Propagator{})

	// Those options are required to properly extract the AWS traceID from the trigger event
	lambda.StartWithOptions(otellambda.InstrumentHandler(Handle, xrayconfig.WithRecommendedOptions(tp)...))
}
