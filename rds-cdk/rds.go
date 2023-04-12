package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type instanceConfig struct {
	DatabaseName string
	AppName      string
	Namespace    string
	VpcId        string
}

// Should move this to read a JSON file
func newInstanceConfig() *instanceConfig {
	db := strings.ReplaceAll(getEnvWithDefault("DATABASE_NAME", "instance"), "-", "")
	app := strings.ReplaceAll(getEnvWithDefault("ACORN_APP", "app"), "-", "")
	ns := strings.ReplaceAll(getEnvWithDefault("ACORN_NAMESPACE", "acorn"), "-", "")
	vpcId := getEnvWithDefault("VPC_ID", "")

	return &instanceConfig{
		DatabaseName: db,
		AppName:      app,
		Namespace:    ns,
		VpcId:        vpcId,
	}
}

func getEnvWithDefault(v, def string) string {
	val := os.Getenv(v)
	if val != "" {
		return val
	}
	return def
}

func (ic *instanceConfig) getQualifiedName(item string) *string {
	c := cases.Title(language.AmericanEnglish)
	return jsii.String(fmt.Sprintf("%s%s%s%s", c.String(ic.AppName), c.String(ic.Namespace), c.String(ic.DatabaseName), c.String(item)))
}

func NewRDSStack(scope constructs.Construct, props *awscdk.StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = *props
	}

	cfg := newInstanceConfig()

	stack := awscdk.NewStack(scope, cfg.getQualifiedName("Stack"), &sprops)

	vpc := awsec2.Vpc_FromLookup(stack, jsii.String("VPC"), &awsec2.VpcLookupOptions{
		VpcId: jsii.String(cfg.VpcId),
	})

	sg := awsec2.NewSecurityGroup(stack, cfg.getQualifiedName("SG"), &awsec2.SecurityGroupProps{
		Vpc:              vpc,
		AllowAllOutbound: jsii.Bool(true),
		Description:      jsii.String("Acorn created Rds security group"),
	})

	subnetGroup := awsrds.NewSubnetGroup(stack, cfg.getQualifiedName("SubnetGroup"), &awsrds.SubnetGroupProps{
		Description: jsii.String("RDS SUBNETS..."),
		Vpc:         vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
	})

	for _, i := range *vpc.PrivateSubnets() {
		sg.AddIngressRule(awsec2.Peer_Ipv4(i.Ipv4CidrBlock()), awsec2.Port_Tcp(jsii.Number(3306)), jsii.String("Allow from private subnets"), jsii.Bool(false))
	}
	for _, i := range *vpc.PublicSubnets() {
		sg.AddIngressRule(awsec2.Peer_Ipv4(i.Ipv4CidrBlock()), awsec2.Port_Tcp(jsii.Number(3306)), jsii.String("Allow from public subnets"), jsii.Bool(false))
	}
	sgs := &[]awsec2.ISecurityGroup{sg}

	creds := awsrds.Credentials_FromGeneratedSecret(jsii.String("clusteradmin"), &awsrds.CredentialsBaseOptions{})

	cluster := awsrds.NewServerlessCluster(stack, cfg.getQualifiedName("Cluster"), &awsrds.ServerlessClusterProps{
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

	awscdk.Tags_Of(cluster).Add(jsii.String("AcornSVCOwner"), cfg.getQualifiedName(""), &awscdk.TagProps{})

	port := "3306"
	pSlice := strings.SplitN(*cluster.ClusterEndpoint().SocketAddress(), ":", 2)
	if len(pSlice) == 2 {
		port = pSlice[1]
	}

	awscdk.NewCfnOutput(stack, cfg.getQualifiedName("host"), &awscdk.CfnOutputProps{
		Value: cluster.ClusterEndpoint().Hostname(),
	})
	awscdk.NewCfnOutput(stack, cfg.getQualifiedName("port"), &awscdk.CfnOutputProps{
		Value: &port,
	})
	awscdk.NewCfnOutput(stack, cfg.getQualifiedName("adminusername"), &awscdk.CfnOutputProps{
		Value: creds.Username(),
	})
	awscdk.NewCfnOutput(stack, cfg.getQualifiedName("adminpasswordarn"), &awscdk.CfnOutputProps{
		Value: cluster.Secret().SecretArn(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	NewRDSStack(app, &awscdk.StackProps{
		Env: rdsenv(),
	})

	app.Synth(nil)
}

func rdsenv() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
