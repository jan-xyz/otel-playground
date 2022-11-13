package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	metricsglobal "go.opentelemetry.io/otel/metric/global"
)

func Handle(ctx context.Context, _ any) (any, error) {
	ctx = context.WithValue(ctx, &reqID{}, "xxxxxxxxxxxxxxxxxxxxxxx")
	ctx, span := otel.Tracer(service).Start(ctx, "Handle")
	defer span.End()

	for i := 0; i < 10; i++ {
		process(ctx, i)
	}

	logrus.WithContext(ctx).
		WithError(errors.New("hello world")).
		WithField("foo", "bar").
		Error("something failed")

	return nil, nil
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
