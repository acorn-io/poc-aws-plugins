package rds

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	configFile = "/app/config.json"
)

var (
	SizeMap = map[string]awsec2.InstanceSize{
		"small":  awsec2.InstanceSize_SMALL,
		"medium": awsec2.InstanceSize_MEDIUM,
		"large":  awsec2.InstanceSize_LARGE,
	}
	acornTags = map[string]string{
		"acorn.io/managed":      "true",
		"acorn.io/project-name": os.Getenv("ACORN_PROJECT"),
		"acorn.io/acorn-name":   os.Getenv("ACORN_NAME"),
		"acorn.io/account-id":   os.Getenv("ACORN_ACCOUNT"),
	}
)

type instanceConfig struct {
	DatabaseName       string            `json:"dbName"`
	InstanceSize       string            `json:"instanceSize"`
	AdminUser          string            `json:"adminUsername"`
	Tags               map[string]string `json:"tags"`
	DeletionProtection bool              `json:"deletionProtection"`
	VpcID              string
}

func instanceConfigFromFile() (*instanceConfig, error) {
	conf := &instanceConfig{}
	fileContent, err := os.ReadFile(configFile)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(fileContent, conf)
	if err != nil {
		return nil, err
	}
	conf.VpcID = os.Getenv("VPC_ID")
	return conf, nil
}

func NewInstanceConfig() (*instanceConfig, error) {
	ic, err := instanceConfigFromFile()
	if err != nil {
		return ic, err
	}

	return ic, nil
}

func AppendGlobalTags(tags *map[string]*string, newTags map[string]string) *map[string]*string {
	result := map[string]*string{}
	if tags != nil {
		for k, v := range *tags {
			result[k] = v
		}
	}
	for k, v := range newTags {
		result[k] = jsii.String(v)
	}
	for k, v := range acornTags {
		if v != "" {
			result[k] = jsii.String(v)
		}
	}
	return &result
}

func GetAllowAllVPCSecurityGroup(scope constructs.Construct, name *string, vpc awsec2.IVpc) awsec2.SecurityGroup {
	sg := awsec2.NewSecurityGroup(scope, name, &awsec2.SecurityGroupProps{
		Vpc:              vpc,
		AllowAllOutbound: jsii.Bool(true),
		Description:      jsii.String("Acorn created RDS security group"),
	})

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
		Description: jsii.String("Acorn created RDS Subnets"),
		Vpc:         vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
	})

	return subnetGroup
}
