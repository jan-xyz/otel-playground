package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
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

	// The code that defines your stack goes here

	_ = awslambda.NewFunction(stack, jsii.String("otel-test"), &awslambda.FunctionProps{
		Code: awslambda.AssetCode_FromAsset(jsii.String("../"), &awss3assets.AssetOptions{
			Bundling: &awscdk.BundlingOptions{
				Local: &bundler{},
				Image: awslambda.Runtime_GO_1_X().BundlingImage(),
				Command: &[]*string{
					jsii.String("go"),
					jsii.String("build"),
					jsii.String("-o"),
					jsii.String("./asset-input"),
					jsii.String("."),
				},
				Environment: &map[string]*string{
					"GOOS":        jsii.String("linux"),
					"GOARCH":      jsii.String("amd64"),
					"CGO_ENABLED": jsii.String("0"),
				},
				User: jsii.String("root"),
			},
		}),
		Handler: jsii.String("otel-playground"),
		Runtime: awslambda.Runtime_GO_1_X(),
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
	cmd := exec.Command("go", "build", "-o", *outputDir, ".")
	var out bytes.Buffer
	cmd.Stderr = &out
	env := os.Environ()
	env = append(env, "GOOS=linux", "GOARCH=amd64")
	cmd.Env = env
	cmd.Dir = ".."
	err := cmd.Run()
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
