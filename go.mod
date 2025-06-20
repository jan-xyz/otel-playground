module otel-playground

go 1.23.0

toolchain go1.24.1

require (
	github.com/aws/aws-lambda-go v1.48.0
	github.com/jan-xyz/box v0.4.1
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda v0.61.0
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig v0.61.0
	go.opentelemetry.io/contrib/propagators/aws v1.36.0
	go.opentelemetry.io/otel v1.36.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.36.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.36.0
	go.opentelemetry.io/otel/sdk v1.36.0
	go.opentelemetry.io/otel/sdk/metric v1.36.0
	go.opentelemetry.io/otel/trace v1.36.0
)

require (
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
)

require (
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	go.opentelemetry.io/contrib/bridges/otelslog v0.11.0
	go.opentelemetry.io/contrib/detectors/aws/lambda v0.61.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.12.2
	go.opentelemetry.io/otel/log v0.12.2
	go.opentelemetry.io/otel/metric v1.36.0
	go.opentelemetry.io/otel/sdk/log v0.12.2
	go.opentelemetry.io/proto/otlp v1.6.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/grpc v1.72.1 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
