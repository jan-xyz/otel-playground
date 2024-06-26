package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jan-xyz/box"
	awslambdago "github.com/jan-xyz/box/transports/github.com/aws/aws-lambda-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

var (
	service     = "otel-test"
	environment = "dev"
	res         = resource.NewWithAttributes(
		semconv.SchemaURL,

		// sets the service correctly
		semconv.ServiceName(service),

		// helps parse stack traces and errors
		semconv.TelemetrySDKLanguageGo,

		// others
		semconv.DeploymentEnvironment(environment),
		semconv.ServiceVersion("v0.0.1"),
	)
)

func main() {
	ctx := context.Background()

	// Setup instrumentation
	lp := setupLogging(ctx, res)
	mp := setupMetrics(ctx, res)
	tp := setupTracing(ctx, res)

	ep := box.Chain(
		tracingMiddleware(),
		loggingMiddleware(),
	)(Endpoint)

	handler := awslambdago.NewAPIGatewayTransport(
		func(_ *events.APIGatewayProxyRequest) (string, error) { return "", nil },
		func(_ string) (*events.APIGatewayProxyResponse, error) { return nil, nil },
		func(err error) *events.APIGatewayProxyResponse { return nil },
		ep,
	)

	defer func() {
		if err := lp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
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
			handler,
			// options to properly extract the AWS traceID from the trigger event,
			// flush the trace provide after each handler event,
			// set the propagator, carrier and trace provider.
			xrayconfig.WithEventToCarrier(),
			xrayconfig.WithPropagator(),
			otellambda.WithTracerProvider(tp),
			otellambda.WithFlusher(lp),
			otellambda.WithFlusher(tp),
			otellambda.WithFlusher(mp),
		),
	)
}
