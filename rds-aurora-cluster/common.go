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

var (
	sizeMap = map[string]awsec2.InstanceSize{
		"small":  awsec2.InstanceSize_SMALL,
		"medium": awsec2.InstanceSize_MEDIUM,
	}
	acornTags = map[string]string{
		"acorn.io/managed": "true",
	}
)

type instanceConfig struct {
	DatabaseName string   `json:"databaseName"`
	AppName      string   `json:"appName"`
	Namespace    string   `json:"namespace"`
	VpcId        string   `json:"vpcId"`
	Public       bool     `json:"public"`
	InstanceSize string   `json:"instanceSize"`
	Users        []string `json:"users"`
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
		Public:       false,
		InstanceSize: "medium",
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

func tagObject(con constructs.Construct) constructs.Construct {
	for k, v := range acornTags {
		awscdk.Tags_Of(con).Add(jsii.String(k), jsii.String(v), &awscdk.TagProps{})
	}
	return con
}

func getAllowAllVPCSecurityGroup(scope constructs.Construct, name *string, vpc awsec2.IVpc) awsec2.SecurityGroup {
	sg := awsec2.NewSecurityGroup(scope, name, &awsec2.SecurityGroupProps{
		Vpc:              vpc,
		AllowAllOutbound: jsii.Bool(true),
		Description:      jsii.String("Acorn created Rds security group"),
	})
	tagObject(sg)

	for _, i := range *vpc.PrivateSubnets() {
		sg.AddIngressRule(awsec2.Peer_Ipv4(i.Ipv4CidrBlock()), awsec2.Port_Tcp(jsii.Number(3306)), jsii.String("Allow from private subnets"), jsii.Bool(false))
	}
	for _, i := range *vpc.PublicSubnets() {
		sg.AddIngressRule(awsec2.Peer_Ipv4(i.Ipv4CidrBlock()), awsec2.Port_Tcp(jsii.Number(3306)), jsii.String("Allow from public subnets"), jsii.Bool(false))
	}
	return sg
}

func getPrivateSubnetGroup(scope constructs.Construct, name *string, vpc awsec2.IVpc) awsrds.SubnetGroup {
	subnetGroup := awsrds.NewSubnetGroup(scope, name, &awsrds.SubnetGroupProps{
		Description: jsii.String("RDS SUBNETS..."),
		Vpc:         vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
	})
	tagObject(subnetGroup)

	return subnetGroup
}
