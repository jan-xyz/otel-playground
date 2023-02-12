package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jan-xyz/box"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	metricsglobal "go.opentelemetry.io/otel/metric/global"
)

var tracingMiddleware = box.Middleware[string, string](func(next box.Endpoint[string, string]) box.Endpoint[string, string] {
	return box.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		ctx = context.WithValue(ctx, &reqID{}, req)
		ctx, span := otel.Tracer(service).Start(ctx, "Handle")
		defer span.End()
		return next(ctx, req)
	})
})

var loggingMiddleware = box.Middleware[string, string](func(next box.Endpoint[string, string]) box.Endpoint[string, string] {
	return box.Endpoint[string, string](func(ctx context.Context, req string) (string, error) {
		resp, err := next(ctx, req)

		logrus.WithContext(ctx).
			WithError(errors.New("hello world")).
			WithField("foo", "bar").
			Error("something failed")
		return resp, err
	})
})

func Endpoint(ctx context.Context, _ string) (string, error) {
	for i := 0; i < 10; i++ {
		process(ctx, i)
	}

	return "", nil
}

func process(ctx context.Context, i int) {
	meter := metricsglobal.Meter("foo")
	c, err := meter.SyncInt64().Counter("request_handled")
	if err != nil {
		log.Printf("counter failed: %s", err)
	}
	ctx, span := otel.Tracer(service).Start(ctx, "proccessing")
	defer span.End()
	logrus.WithContext(ctx).
		WithField("iteration", i).
		Info("tick!")
	c.Add(ctx, 1)
	<-time.After(3 * time.Millisecond)
}
