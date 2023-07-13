package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	configFile = "/app/config.json"
)

var (
	acornTags = map[string]string{
		"acorn.io/managed":      "true",
		"acorn.io/project-name": os.Getenv("ACORN_PROJECT"),
		"acorn.io/acorn-name":   os.Getenv("ACORN_NAME"),
		"acorn.io/account-id":   os.Getenv("ACORN_ACCOUNT"),
	}
)

func ConfigBytes() ([]byte, error) {
	return os.ReadFile(configFile)
}

func NewConfig(props any) error {
	conf, err := ConfigBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(conf, props)
}

func GetAcornTaggedApp(props *awscdk.AppProps) awscdk.App {
	app := awscdk.NewApp(props)
	AppendScopedTags(app, acornTags)
	return app
}

func AppendScopedTags(scope constructs.Construct, tags map[string]string) {
	scopedTags := awscdk.Tags_Of(scope)
	for k, v := range tags {
		scopedTags.Add(jsii.String(k), jsii.String(v), &awscdk.TagProps{})
	}
}
