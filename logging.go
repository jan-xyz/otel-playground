package main

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func buildLogger() {
	otelhook := otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	))
	logrus.AddHook(otelhook)
	lhook := &hook{
		levels: otelhook.Levels(),
		ctxKeys: []fmt.Stringer{
			&reqID{},
			&userID{},
		},
	}
	logrus.AddHook(lhook)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
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
func (h hook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)

	fields := logrus.Fields{
		"trace_id": span.SpanContext().TraceID().String(),
		"span_id":  span.SpanContext().SpanID().String(),
	}

	for _, v := range h.ctxKeys {
		fields[v.String()] = ctx.Value(v)
	}

	if entry.Caller != nil {
		if entry.Caller.Function != "" {
			fields[string(semconv.CodeFunctionKey)] = entry.Caller.Function
		}
		if entry.Caller.File != "" {
			fields[string(semconv.CodeFilepathKey)] = entry.Caller.File
			fields[string(semconv.CodeLineNumberKey)] = entry.Caller.Line
		}
	}
	for k, v := range entry.Data {
		if k == "error" {
			if err, ok := v.(error); ok {
				typ := reflect.TypeOf(err).String()
				fields[string(semconv.ExceptionTypeKey)] = typ
				fields[string(semconv.ExceptionMessageKey)] = err.Error()
				continue
			}
		}

		fields[k] = v
	}
	*entry = *entry.WithFields(fields)

	return nil
}
