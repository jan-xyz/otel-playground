package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jan-xyz/box"
	awslambdago "github.com/jan-xyz/box/handler/github.com/aws/aws-lambda-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
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
	mp := setupMetrics(ctx, res)
	tp := setupTracing(ctx, res)

	ep := box.Chain(
		tracingMiddleware(),
		loggingMiddleware(),
	)(Endpoint)

	handler := awslambdago.NewAPIGatewayHandler(
		func(_ *events.APIGatewayProxyRequest) (string, error) { return "", nil },
		func(_ string) (*events.APIGatewayProxyResponse, error) { return nil, nil },
		func(err error) (*events.APIGatewayProxyResponse, error) { return nil, nil },
		ep,
	)

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	// Setup Lambda
	lambda.StartWithOptions(
		// wraps the lambda handler to inject the OTEL context
		otellambda.InstrumentHandler(
			handler.Handle,
			// options to properly extract the AWS traceID from the trigger event,
			// flush the trace provide after each handler event,
			// set the propagator, carrier and trace provider.
			xrayconfig.WithEventToCarrier(),
			xrayconfig.WithPropagator(),
			otellambda.WithTracerProvider(tp),
			otellambda.WithFlusher(tp),
			otellambda.WithFlusher(mp),
		),
	)
}
