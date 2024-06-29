package main

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jan-xyz/box"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	logger = otelslog.NewLogger(service)
	tracer = otel.Tracer(service)
	meter  = otel.Meter(service)
)

var (
	requestCounter metric.Int64Counter
	errorCounter   metric.Int64Counter
)

func init() {
	var err error
	requestCounter, err = meter.Int64Counter("request.counter",
		metric.WithDescription("The number of requests"),
		metric.WithUnit("{request}"))
	if err != nil {
		panic(err)
	}

	errorCounter, err = meter.Int64Counter("error.counter",
		metric.WithDescription("The number of errors"),
		metric.WithUnit("{error}"))
	if err != nil {
		panic(err)
	}
}

func metricMiddleware() box.Middleware[string, string] {
	return func(next box.Endpoint[string, string]) box.Endpoint[string, string] {
		return func(ctx context.Context, req string) (string, error) {
			requestCounter.Add(ctx, 1)
			out, err := next(ctx, req)
			if err != nil {
				errorCounter.Add(ctx, 1)
			}
			return out, err
		}
	}
}

func tracingMiddleware() box.Middleware[string, string] {
	return func(next box.Endpoint[string, string]) box.Endpoint[string, string] {
		return func(ctx context.Context, req string) (string, error) {
			ctx = context.WithValue(ctx, &reqID{}, req)
			ctx, span := tracer.Start(ctx, "Handle")
			defer span.End()
			return next(ctx, req)
		}
	}
}

func loggingMiddleware() box.Middleware[string, string] {
	return func(next box.Endpoint[string, string]) box.Endpoint[string, string] {
		return func(ctx context.Context, req string) (string, error) {
			resp, err := next(ctx, req)

			logger.ErrorContext(ctx,
				"something failed",
				slog.String("error", errors.New("hello world").Error()),
				slog.String("foo", "bar"))
			return resp, err
		}
	}
}

func Endpoint(ctx context.Context, _ string) (string, error) {
	for i := 0; i < 10; i++ {
		process(ctx, i)
	}
	return "", nil
}

func process(ctx context.Context, i int) {
	ctx, span := tracer.Start(ctx, "proccessing")
	defer span.End()
	logger.InfoContext(ctx, "tick!", slog.Int("iteration", i))
	<-time.After(3 * time.Millisecond)
}
