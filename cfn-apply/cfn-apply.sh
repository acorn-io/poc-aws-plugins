#!/bin/sh

file="/acorn/data/cfn.yaml"
templateName=${1}

echo "${file} found.."

aws cloudformation deploy --template-file /acorn/data/cfn.yaml --stack-name ${templateName} --capabilities CAPABILITY_IAM --capabilities CAPABILITY_NAMED_IAM --no-cli-pager
