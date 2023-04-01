module otel-playground

go 1.19

require (
	github.com/aws/aws-lambda-go v1.37.0
	github.com/jan-xyz/box v0.1.0
	github.com/sirupsen/logrus v1.9.0
	github.com/uptrace/opentelemetry-go-extra/otellogrus v0.1.17
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda v0.37.0
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig v0.37.0
	go.opentelemetry.io/contrib/propagators/aws v1.15.0
	go.opentelemetry.io/otel v1.14.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.37.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.11.2
	go.opentelemetry.io/otel/metric v0.37.0
	go.opentelemetry.io/otel/sdk v1.14.0
	go.opentelemetry.io/otel/sdk/metric v0.37.0
	go.opentelemetry.io/otel/trace v1.14.0
	google.golang.org/grpc v1.54.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.13.0 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelutil v0.1.17 // indirect
	go.opentelemetry.io/contrib/detectors/aws/lambda v0.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.14.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.11.2 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
