package main

import (
	"os"
	"strings"

	"github.com/acorn-io/aws/rds"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/sirupsen/logrus"
)

func NewRDSStack(scope constructs.Construct, props *awscdk.StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = *props
	}

	cfg, err := rds.NewInstanceConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	sprops.Tags = rds.AppendGlobalTags(sprops.Tags, cfg.Tags)

	stack := awscdk.NewStack(scope, jsii.String("Stack"), &sprops)

	vpc := awsec2.Vpc_FromLookup(stack, jsii.String("VPC"), &awsec2.VpcLookupOptions{
		VpcId: jsii.String(cfg.VpcID),
	})

	subnetGroup := rds.GetPrivateSubnetGroup(stack, jsii.String("SubnetGroup"), vpc)
	sgs := &[]awsec2.ISecurityGroup{
		rds.GetAllowAllVPCSecurityGroup(stack, jsii.String("SG"), vpc),
	}

	creds := awsrds.Credentials_FromGeneratedSecret(jsii.String(cfg.AdminUser), &awsrds.CredentialsBaseOptions{})

	cluster := awsrds.NewServerlessCluster(stack, jsii.String("Cluster"), &awsrds.ServerlessClusterProps{
		Engine:              awsrds.DatabaseClusterEngine_AURORA_MYSQL(),
		DefaultDatabaseName: jsii.String(cfg.DatabaseName),
		CopyTagsToSnapshot:  jsii.Bool(true),
		Credentials:         creds,
		Vpc:                 vpc,
		Scaling: &awsrds.ServerlessScalingOptions{
			AutoPause: awscdk.Duration_Minutes(jsii.Number(10)),
		},
		SubnetGroup:    subnetGroup,
		SecurityGroups: sgs,
	})

	port := "3306"
	pSlice := strings.SplitN(*cluster.ClusterEndpoint().SocketAddress(), ":", 2)
	if len(pSlice) == 2 {
		port = pSlice[1]
	}

	awscdk.NewCfnOutput(stack, jsii.String("host"), &awscdk.CfnOutputProps{
		Value: cluster.ClusterEndpoint().Hostname(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("port"), &awscdk.CfnOutputProps{
		Value: &port,
	})
	awscdk.NewCfnOutput(stack, jsii.String("adminusername"), &awscdk.CfnOutputProps{
		Value: creds.Username(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("adminpasswordarn"), &awscdk.CfnOutputProps{
		Value: cluster.Secret().SecretArn(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	NewRDSStack(app, &awscdk.StackProps{
		Synthesizer: awscdk.NewDefaultStackSynthesizer(&awscdk.DefaultStackSynthesizerProps{
			GenerateBootstrapVersionRule: jsii.Bool(false),
		}),
		Env: &awscdk.Environment{
			Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
			Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
		},
	})

	app.Synth(nil)
}
