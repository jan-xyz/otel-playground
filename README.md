# otel-playground

This is a playground to play with AWS and OTEL in Go lambdas.

## Resources

* [OTLP setup](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp#example-NewExporter)
* [ADOT Lambda setup](https://aws-otel.github.io/docs/getting-started/lambda/lambda-go)
* [ADOT collector default config](https://github.com/aws-observability/aws-otel-collector/blob/main/config.yaml)
* [X-Ray segment documentation](https://docs.aws.amazon.com/xray/latest/devguide/xray-api-segmentdocuments.html)
* [OTEL span to XRay segment conversion](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/exporter/awsxrayexporter/internal/translator/segment.go)
* [OTEL error to XRay cause conversion](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/exporter/awsxrayexporter/internal/translator/cause.go)
