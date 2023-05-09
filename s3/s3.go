package main

import (
	"os"
	"strconv"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type S3StackProps struct {
	awscdk.StackProps
	S3BucketName            string
	S3Versioned             bool
	S3DeleteBucketOnDestroy bool
}

func NewS3Stack(scope constructs.Construct, id string, props *S3StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	removalPolicy := awscdk.RemovalPolicy_RETAIN
	if props.S3DeleteBucketOnDestroy == true {
		removalPolicy = awscdk.RemovalPolicy_DESTROY
	}

	s3Bucket := awss3.NewBucket(stack, jsii.String(props.S3BucketName), &awss3.BucketProps{
		Versioned:         jsii.Bool(props.S3Versioned),
		BucketName:        jsii.String(props.S3BucketName),
		AutoDeleteObjects: jsii.Bool(props.S3DeleteBucketOnDestroy),
		RemovalPolicy:     removalPolicy,
	})

	awscdk.NewCfnOutput(stack, jsii.String("s3BucketArn"), &awscdk.CfnOutputProps{
		Value: s3Bucket.BucketArn(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	s3BucketName := os.Getenv("S3_BUCKET_NAME")
	s3BucketVersioned, err := strconv.ParseBool(os.Getenv("S3_BUCKET_VERSIONED"))
	s3DelteBucketOnDestroy, err := strconv.ParseBool(os.Getenv("S3_DELETE_BUCKET_ON_DESTROY"))

	if err != nil {
		os.Exit(1)
	}

	S3StackName := "S3-" + s3BucketName

	NewS3Stack(app, S3StackName, &S3StackProps{
		StackProps: awscdk.StackProps{
			Synthesizer: awscdk.NewDefaultStackSynthesizer(&awscdk.DefaultStackSynthesizerProps{
				GenerateBootstrapVersionRule: jsii.Bool(false),
			}),
			Env: env(),
		},
		S3BucketName:            s3BucketName,
		S3Versioned:             s3BucketVersioned,
		S3DeleteBucketOnDestroy: s3DelteBucketOnDestroy,
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
