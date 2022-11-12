package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	metricsglobal "go.opentelemetry.io/otel/metric/global"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	service           = "otel-test"
	environment       = "dev"
	id          int64 = 12345
)

type reqID struct{}

func Handle(ctx context.Context) error {
	ctx = context.WithValue(ctx, reqID{}, "xxxxxxxxxxxxxxxxxxxxxxx")
	ctx, span := otel.Tracer(service).Start(ctx, "Handle")
	defer span.End()

	meter := metricsglobal.Meter("foo")
	c, err := meter.SyncInt64().Counter("request_handled")
	if err != nil {
		log.Printf("counter failed: %s", err)
	}
	for i := 0; i < 10; i++ {
		c.Add(ctx, 1)
	}

	logrus.WithContext(ctx).
		WithError(errors.New("hello world")).
		WithField("foo", "bar").
		Error("something failed")

	return nil
}

func buildProvider() (*tracesdk.TracerProvider, *metricsdk.MeterProvider, error) {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(service),
		attribute.String("environment", environment),
		attribute.Int64("ID", id),
	)

	// Create Metrics exporter
	m_exp, err := stdoutmetric.New()
	if err != nil {
		return nil, nil, fmt.Errorf("creating metric exporter: %w", err)
	}
	mp := metricsdk.NewMeterProvider(
		metricsdk.WithReader(metricsdk.NewPeriodicReader(m_exp)),
		// Record information about this application in a Resource.
		metricsdk.WithResource(res),
	)

	// Create Trace exporter
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, nil, fmt.Errorf("creating tace exporter: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(res),
	)

	otelhook := otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	))
	logrus.AddHook(otelhook)
	logrus.AddHook(&hook{levels: otelhook.Levels()})
	return tp, mp, nil
}

func main() {
	ctx := context.Background()
	tp, mp, err := buildProvider()
	if err != nil {
		panic(err)
	}
	defer mp.ForceFlush(ctx)
	defer tp.ForceFlush(ctx)
	otel.SetTracerProvider(tp)
	metricsglobal.SetMeterProvider(mp)

	Handle(ctx)
}

type hook struct {
	levels []logrus.Level
}

// Levels returns logrus levels on which this hook is fired.
func (h hook) Levels() []logrus.Level {
	return h.levels
}

func (*hook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)

	*entry = *entry.WithFields(logrus.Fields{
		"trace_id": span.SpanContext().TraceID().String(),
		"span_id":  span.SpanContext().SpanID().String(),
		"req_id":   ctx.Value(reqID{}),
	})

	return nil
}
