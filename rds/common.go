package rds

import (
	"encoding/json"
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
	SizeMap = map[string]awsec2.InstanceSize{
		"small":  awsec2.InstanceSize_SMALL,
		"medium": awsec2.InstanceSize_MEDIUM,
	}
	configFile = "/app/config.json"
	acornTags  = map[string]string{
		"acorn.io/managed": "true",
	}
	tags = map[string]string{}
)

type instanceConfig struct {
	DatabaseName string            `json:"databaseName"`
	AppName      string            `json:"appName"`
	Namespace    string            `json:"namespace"`
	VpcId        string            `json:"vpcId"`
	Public       bool              `json:"public"`
	InstanceSize string            `json:"instanceSize"`
	AdminUser    string            `json:"adminUser"`
	Tags         map[string]string `json:"tags"`
}

func instanceConfigFromFile() (*instanceConfig, error) {
	conf := &instanceConfig{}
	fileContent, err := os.ReadFile(configFile)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(fileContent, conf)
	return conf, err
}

// Should move this to read a JSON file
func NewInstanceConfig() (*instanceConfig, error) {
	ic, err := instanceConfigFromFile()
	if err != nil {
		return ic, err
	}

	ic.AppName = getEnvWithDefault("ACORN_APP", "app")
	ic.Namespace = getEnvWithDefault("ACORN_NAMESPACE", "acorn")
	ic.VpcId = getEnvWithDefault("VPC_ID", "")

	setTags(ic.Tags)

	return ic, nil
}

func getEnvWithDefault(v, def string) string {
	val := os.Getenv(v)
	if val != "" {
		return val
	}
	return def
}

func (ic *instanceConfig) GetQualifiedName(item string) *string {
	db := strings.ReplaceAll(ic.DatabaseName, "-", "")
	app := strings.ReplaceAll(ic.AppName, "-", "")
	ns := strings.ReplaceAll(ic.Namespace, "-", "")

	c := cases.Title(language.AmericanEnglish)
	return jsii.String(fmt.Sprintf("%s%s%s%s", c.String(app), c.String(ns), c.String(db), c.String(item)))
}

func setTags(t map[string]string) {
	for k, v := range t {
		tags[k] = v
	}
	for k, v := range acornTags {
		tags[k] = v
	}
}

func TagObject(con constructs.Construct) constructs.Construct {
	for k, v := range tags {
		awscdk.Tags_Of(con).Add(jsii.String(k), jsii.String(v), &awscdk.TagProps{})
	}
	return con
}

func GetAllowAllVPCSecurityGroup(scope constructs.Construct, name *string, vpc awsec2.IVpc) awsec2.SecurityGroup {
	sg := awsec2.NewSecurityGroup(scope, name, &awsec2.SecurityGroupProps{
		Vpc:              vpc,
		AllowAllOutbound: jsii.Bool(true),
		Description:      jsii.String("Acorn created Rds security group"),
	})
	TagObject(sg)

	for _, i := range *vpc.PrivateSubnets() {
		sg.AddIngressRule(awsec2.Peer_Ipv4(i.Ipv4CidrBlock()), awsec2.Port_Tcp(jsii.Number(3306)), jsii.String("Allow from private subnets"), jsii.Bool(false))
	}
	for _, i := range *vpc.PublicSubnets() {
		sg.AddIngressRule(awsec2.Peer_Ipv4(i.Ipv4CidrBlock()), awsec2.Port_Tcp(jsii.Number(3306)), jsii.String("Allow from public subnets"), jsii.Bool(false))
	}
	return sg
}

func GetPrivateSubnetGroup(scope constructs.Construct, name *string, vpc awsec2.IVpc) awsrds.SubnetGroup {
	subnetGroup := awsrds.NewSubnetGroup(scope, name, &awsrds.SubnetGroupProps{
		Description: jsii.String("RDS SUBNETS..."),
		Vpc:         vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
	})
	TagObject(subnetGroup)

	return subnetGroup
}
