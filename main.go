package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
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
	tp, mp, err := buildProvider()
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)
	metricsglobal.SetMeterProvider(mp)

	handle := middlewareChain(flushOtel(mp, tp))(Handle)

	lambda.StartWithOptions(handle, lambda.WithContext(ctx))
}
