#!/bin/bash

stack="${1}"

aws cloudformation describe-stacks --stack-name ${stack} --query 'Stacks[0].Outputs' > outputs.json

port="$(jq -r '.[] | select(.OutputKey=="AcornRdsClusterport")|.OutputValue' outputs.json )"
address="$(jq -r '.[] | select(.OutputKey=="AcornRdsClusterhost")|.OutputValue' outputs.json )"
admin_username="$(jq -r '.[] | select(.OutputKey=="AcornRdsClusterusername")|.OutputValue' outputs.json )"
password_arn="$(jq -r '.[] | select(.OutputKey=="AcornRdsClusterpasswordarn")|.OutputValue' outputs.json )"
admin_password="$(aws --output json secretsmanager get-secret-value --secret-id $password_arn --query 'SecretString' | jq -r .|jq -r .password)"

cat > /run/secrets/output<<EOF
services: db: {
    default: true
    address: "${address}"
    secrets: ["admin"]
    ports: [
        {
            port: ${port}
        }
    ]
    data: dbname: "instance"
}

secrets: "admin": {
	type: "basic"
	data: {
        username: "${admin_username}"
        password: "${admin_password}"
    }
}
EOF
