package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	metricsglobal "go.opentelemetry.io/otel/metric/global"
)

var (
	service           = "otel-test"
	environment       = "dev"
	id          int64 = 12345
)

func main() {
	ctx := context.Background()
	tp, mp, err := buildProvider(ctx)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)
	metricsglobal.SetMeterProvider(mp)
	defer tp.Shutdown(ctx)
	defer mp.Shutdown(ctx)

	otel.SetTextMapPropagator(xray.Propagator{})

	lambda.StartWithOptions(otellambda.InstrumentHandler(Handle, xrayconfig.WithRecommendedOptions(tp)...))
}
