#!/bin/sh

set -e

file="/acorn/data/cfn.yaml"

stackName="${1}"
if [ $# -eq 0 ]; then
    echo "No stack name defined..."
    echo "syntax:  $0 [STACKNAME]"
    exit 1
fi 

sanitizedName=$(sed 's/-//g' <<< ${stackName})

echo "${file} found.."

aws cloudformation deploy --template-file ${file} --stack-name "${sanitizedName}" --capabilities CAPABILITY_IAM --capabilities CAPABILITY_NAMED_IAM --no-cli-pager
aws cloudformation describe-stacks --stack-name "${sanitizedName}" --query 'Stacks[0].Outputs' > outputs.json

if [ -f /app/scripts/render.sh ]; then
  /app/scripts/render.sh
fi
