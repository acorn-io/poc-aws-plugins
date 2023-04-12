#!/bin/sh

set -e

file="/acorn/data/cfn.yaml"
APP=$(sed 's/-//g' <<< ${ACORN_APP^})
NS=$(sed 's/-//g' <<< ${ACORN_NAMESPACE^})
DB=$(sed 's/-//g' <<< ${DATABASE_NAME^})

stackName="${APP}${NS}${DB}Stack"


echo "${file} found.."

aws cloudformation deploy --template-file /acorn/data/cfn.yaml --stack-name ${stackName} --capabilities CAPABILITY_IAM --capabilities CAPABILITY_NAMED_IAM --no-cli-pager
aws cloudformation describe-stacks --stack-name ${stackName} --query 'Stacks[0].Outputs' > outputs.json

if [ -f /app/render.sh ]; then
  /app/render.sh
fi