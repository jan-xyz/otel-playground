package main

import (
	"context"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	endpoint   = func(context.Context, any) (any, error)
	middleware = func(endpoint) endpoint
)

func middlewareChain(outer middleware, others ...middleware) middleware {
	return func(next endpoint) endpoint {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

func flushOtel(mp *metric.MeterProvider, tp *trace.TracerProvider) middleware {
	return func(next endpoint) endpoint {
		return func(ctx context.Context, req any) (any, error) {
			defer mp.ForceFlush(ctx)
			defer tp.ForceFlush(ctx)

			return next(ctx, req)
		}
	}
}
