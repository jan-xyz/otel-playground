package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"

	"go.opentelemetry.io/otel/trace"
)

func setupLogging() *log.LoggerProvider {
	exporter, err := stdoutlog.New()
	if err != nil {
		panic(err)
	}

	provider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
		log.WithResource(res),
	)
	global.SetLoggerProvider(provider)
	return provider
}

type userID struct{}

func (userID) String() string {
	return "user_id"
}

type reqID struct{}

func (reqID) String() string {
	return "req_id"
}

type ContextHandler struct {
	slog.Handler
	ctxKeys []fmt.Stringer
}

// Handle adds contextual attributes to the Record before calling the underlying
// handler
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)

	r.AddAttrs(slog.String("trace_id", span.SpanContext().TraceID().String()))
	r.AddAttrs(slog.String("span_id", span.SpanContext().SpanID().String()))

	// required for correlating logs in x-ray to the correct span
	r.AddAttrs(slog.String("_X_AMAZN_TRACE_ID", os.Getenv("_X_AMZN_TRACE_ID")))

	// adding custom fields from the context to the logs
	for _, v := range h.ctxKeys {
		if value := ctx.Value(v); value != nil {
			r.AddAttrs(slog.Any(v.String(), value))
		}
	}

	return h.Handler.Handle(ctx, r)
}
