package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	otelCollectorVersion = "0-90-1"
	architecture         = "arm64"
)

type InfraStackProps struct {
	awscdk.StackProps
}

func NewInfraStack(scope constructs.Construct, id string, props *InfraStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	adotLayer := fmt.Sprintf("arn:aws:lambda:%s:901920570463:layer:aws-otel-collector-%s-ver-%s:1", *stack.Region(), architecture, otelCollectorVersion)
	lambdaArch := awslambda.Architecture_X86_64()
	if architecture == "arm64" {
		lambdaArch = awslambda.Architecture_ARM_64()
	}

	_ = awslambda.NewFunction(stack, jsii.String("otel-test"), &awslambda.FunctionProps{
		Code: awslambda.AssetCode_FromAsset(jsii.String("../"), &awss3assets.AssetOptions{
			Bundling: &awscdk.BundlingOptions{
				Local: &bundler{},
				Image: awslambda.Runtime_PROVIDED_AL2().BundlingImage(),
			},
		}),
		Handler:      jsii.String("bootstrap"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Architecture: lambdaArch,
		Layers: &[]awslambda.ILayerVersion{awslambda.LayerVersion_FromLayerVersionArn(
			stack,
			jsii.String("adot"),
			// layer arn: arn:aws:lambda:<region>:901920570463:layer:aws-otel-collector-<architecture>-ver-0-66-0:1
			&adotLayer,
		)},
		Tracing: awslambda.Tracing_ACTIVE,
		Environment: &map[string]*string{
			"OPENTELEMETRY_COLLECTOR_CONFIG_FILE": jsii.String("/var/task/collector.yaml"),
		},
	})
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewInfraStack(app, "InfraStack", &InfraStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

type bundler struct{}

func (bundler) TryBundle(outputDir *string, options *awscdk.BundlingOptions) *bool {
	// bundle binary
	cmd := exec.Command("cp", "collector.yaml", *outputDir)
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Dir = ".."
	err := cmd.Run()
	if err != nil {
		fmt.Println("failed to copy:", err, out.String())
		return jsii.Bool(false)

	}
	// bundle collector.yaml
	cmd = exec.Command("go", "build", "-o", path.Join(*outputDir, "bootstrap"), ".")
	cmd.Stderr = &out
	env := os.Environ()
	env = append(env, "GOOS=linux", fmt.Sprintf("GOARCH=%s", architecture))
	cmd.Env = env
	cmd.Dir = ".."
	err = cmd.Run()
	if err != nil {
		fmt.Println("failed to build:", err, out.String())
		return jsii.Bool(false)

	}
	return jsii.Bool(true)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
