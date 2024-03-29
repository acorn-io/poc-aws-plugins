package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/sirupsen/logrus"
)

type MyStackProps struct {
	awscdk.StackProps
	MakePublic bool              `json:"makePublic" yaml:"makePublic"`
	Versioned  bool              `json:"versioned" yaml:"versioned"`
	BucketName string            `json:"bucketName" yaml:"bucketName"`
	UserTags   map[string]string `json:"tags" yaml:"tags"`
}

func NewMyStack(scope constructs.Construct, id string, props *MyStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, jsii.String(id), &sprops)

	// Create an S3 bucket
	bucket := awss3.NewBucket(stack, jsii.String(props.BucketName), &awss3.BucketProps{
		PublicReadAccess: jsii.Bool(props.MakePublic),
		Versioned:        jsii.Bool(props.Versioned),
		RemovalPolicy:    awscdk.RemovalPolicy_DESTROY,
	})

	// Output the bucket URL and ARN
	awscdk.NewCfnOutput(stack, jsii.String("BucketURL"), &awscdk.CfnOutputProps{
		Value: bucket.BucketWebsiteUrl(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("BucketARN"), &awscdk.CfnOutputProps{
		Value: bucket.BucketArn(),
	})

	return stack
}

func main() {
	app := GetAcornTaggedApp(nil)

	stackProps := &MyStackProps{
		StackProps: awscdk.StackProps{
			Synthesizer: awscdk.NewDefaultStackSynthesizer(&awscdk.DefaultStackSynthesizerProps{
				GenerateBootstrapVersionRule: jsii.Bool(false),
			}),
			Env: &awscdk.Environment{
				Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
				Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
			},
		},
	}

	// Read config from file
	if err := NewConfig(stackProps); err != nil {
		logrus.Fatal(err)
	}

	AppendScopedTags(app, stackProps.UserTags)

	NewMyStack(app, "MyStack", stackProps)

	app.Synth(nil)
}
