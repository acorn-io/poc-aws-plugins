#!/bin/sh

file="/acorn/data/cfn.yaml"
templateName=${1}

echo "waiting for file ${file}..."
while [ ! -f "$file" ];do
  echo "${file} still not present..."
  sleep 1
done

echo "${file} found.."
sleep 1

aws cloudformation deploy --template-file /acorn/data/cfn.yaml --stack-name ${templateName} --capabilities CAPABILITY_IAM --capabilities CAPABILITY_NAMED_IAM --no-cli-pager
