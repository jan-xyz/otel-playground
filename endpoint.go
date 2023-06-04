package main

import (
	"context"
	"errors"
	"time"

	"github.com/jan-xyz/box"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

var tracingMiddleware = func() box.Middleware[string, string] {
	return func(next box.Endpoint[string, string]) box.Endpoint[string, string] {
		return func(ctx context.Context, req string) (string, error) {
			ctx = context.WithValue(ctx, &reqID{}, req)
			ctx, span := otel.Tracer(service).Start(ctx, "Handle")
			defer span.End()
			return next(ctx, req)
		}
	}
}

var loggingMiddleware = func() box.Middleware[string, string] {
	return func(next box.Endpoint[string, string]) box.Endpoint[string, string] {
		return func(ctx context.Context, req string) (string, error) {
			resp, err := next(ctx, req)

			logrus.WithContext(ctx).
				WithError(errors.New("hello world")).
				WithField("foo", "bar").
				Error("something failed")
			return resp, err
		}
	}
}

func Endpoint(ctx context.Context, _ string) (string, error) {
	meter := otel.Meter("")
	for i := 0; i < 10; i++ {
		process(ctx, i)
	}
	if c, err := meter.Int64Counter("request_handled"); err == nil {
		c.Add(ctx, 1)
	}
	return "", nil
}

func process(ctx context.Context, i int) {
	ctx, span := otel.Tracer(service).Start(ctx, "proccessing")
	defer span.End()
	logrus.WithContext(ctx).
		WithField("iteration", i).
		Info("tick!")
	<-time.After(3 * time.Millisecond)
}
