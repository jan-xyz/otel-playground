package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	service           = "otel-test"
	environment       = "dev"
	id          int64 = 12345
	res               = resource.NewWithAttributes(
		semconv.SchemaURL,
		// sets the service correctly
		semconv.ServiceNameKey.String(service),

		// helps parse stack traces and errors
		semconv.TelemetrySDKLanguageGo,

		// others
		semconv.DeploymentEnvironmentKey.String(environment),
		attribute.Int64("id", id),
	)
)

func main() {
	ctx := context.Background()

	// Setup instrumentation
	setupLogging()
	setupMetrics(ctx, res)
	tp := setupTracing(ctx, res)

	// Setup Lambda
	lambda.StartWithOptions(
		// wraps the lambda handler to inject the OTEL context
		otellambda.InstrumentHandler(
			Handle,
			// options to properly extract the AWS traceID from the trigger event,
			// flush the trace provide after each handler event,
			// set the propagator, carrier and trace provider.
			xrayconfig.WithRecommendedOptions(tp)...,
		),
	)
}
