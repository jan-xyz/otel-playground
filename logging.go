package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func buildLogger() {
	// hook to add logs to the span
	otelhook := otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	))
	logrus.AddHook(otelhook)

	// hook to add tracing and context information to the logs
	lhook := &hook{
		levels: otelhook.Levels(),
		// these are context keys that will be added to each log-line
		ctxKeys: []fmt.Stringer{
			&reqID{},
			&userID{},
		},
	}
	logrus.AddHook(lhook)

	// add caller information (Function, File, Line no) to logging entries
	logrus.SetReportCaller(true)

	// set to use JSON logging for easier parsing
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

type userID struct{}

func (userID) String() string {
	return "user_id"
}

type reqID struct{}

func (reqID) String() string {
	return "req_id"
}

type hook struct {
	levels  []logrus.Level
	ctxKeys []fmt.Stringer
}

// Levels returns logrus levels on which this hook is fired.
func (h hook) Levels() []logrus.Level {
	return h.levels
}

// took inspiration from https://github.com/uptrace/uptrace-go/blob/extra/otellogrus/v1.1.0/extra/otellogrus/otellogrus.go
func (h hook) Fire(e *logrus.Entry) error {
	ctx := e.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)

	e.Data["trace_id"] = span.SpanContext().TraceID().String()

	e.Data["span_id"] = span.SpanContext().SpanID().String()

	// required for correlating logs in x-ray to the correct span
	e.Data["_X_AMAZN_TRACE_ID"] = os.Getenv("_X_AMZN_TRACE_ID")

	// adding custom fields from the context to the logs
	for _, v := range h.ctxKeys {
		if value := ctx.Value(v); value != nil {
			e.Data[v.String()] = value
		}
	}

	// add caller informaiton. To use thise you need to call
	// 	logrus.SetReportCaller(true)
	// when setting up the logger
	if e.Caller != nil {
		if e.Caller.Function != "" {
			e.Data[string(semconv.CodeFunctionKey)] = e.Caller.Function
		}
		if e.Caller.File != "" {
			e.Data[string(semconv.CodeFilepathKey)] = e.Caller.File
			e.Data[string(semconv.CodeLineNumberKey)] = e.Caller.Line
		}
	}

	return nil
}
