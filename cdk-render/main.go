package main

import (
	"os"

	"github.com/acorn-io/aws-acorn/pkg/config"
	ctx "github.com/acorn-io/aws-acorn/pkg/context"
	"github.com/sirupsen/logrus"
)

func main() {
	awsConfigPath := "/acorn/aws"
	if p := os.Getenv("AWS_CONFIG_PATH"); p != "" {
		awsConfigPath = p
	}

	output := "./cdk.context.json"
	if od := os.Getenv("OUTPUT_FILE"); od != "" {
		output = od
	}
	logrus.Infof("Reading config dir %s", awsConfigPath)
	cfg, err := config.ReadData(awsConfigPath)
	if err != nil {
		logrus.Fatal(err)
	}

	acct, ok := cfg["account-id"]
	if !ok {
		logrus.Fatal("AWS account id not found in file `account-id`")
	}

	reg, ok := cfg["aws-region"]
	if !ok {
		logrus.Fatal("AWS region not found in file `region`")
	}

	cdkContext, err := ctx.NewContext(acct, reg)
	if err != nil {
		logrus.Fatal(err)
	}

	vpcId, ok := cfg["vpc-id"]
	if !ok {
		logrus.Fatal("VPC ID not found in file `vpc-id`")
	}
	vpcPlugin := ctx.NewVpcPlugin(vpcId)
	cdkContext.AddPlugin(vpcPlugin)

	if err := ctx.WriteFile(output, cdkContext); err != nil {
		logrus.Fatal(err)
	}
}
