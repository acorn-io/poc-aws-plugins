package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsaps"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AmpStackProps struct {
	StackProps    awscdk.StackProps
	WorkspaceName string
}

func NewAmpStack(scope constructs.Construct, id string, props *AmpStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	ampWorkspace := awsaps.NewCfnWorkspace(stack, jsii.String(props.WorkspaceName), &awsaps.CfnWorkspaceProps{
		Alias: jsii.String(props.WorkspaceName),
	})

	awscdk.NewCfnOutput(stack, jsii.String("AMPEndpointURL"), &awscdk.CfnOutputProps{
		Value: ampWorkspace.AttrPrometheusEndpoint(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("AMPWorkspaceArn"), &awscdk.CfnOutputProps{
		Value: ampWorkspace.AttrArn(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	name := os.Getenv("WORKSPACE_NAME")

	NewAmpStack(app, "AmpStack", &AmpStackProps{
		StackProps: awscdk.StackProps{
			Synthesizer: awscdk.NewDefaultStackSynthesizer(&awscdk.DefaultStackSynthesizerProps{
				GenerateBootstrapVersionRule: jsii.Bool(false),
			}),
			Env: env(),
		},
		WorkspaceName: name,
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
