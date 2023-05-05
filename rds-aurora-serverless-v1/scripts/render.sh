#!/bin/bash

set -e 

APP=$(sed 's/-//g' <<< ${ACORN_APP^})
NS=$(sed 's/-//g' <<< ${ACORN_NAMESPACE^})
DB=$(sed 's/-//g' <<< ${DATABASE_NAME^})

stackPfx="${APP}${NS}${DB}"


port="$(jq -r --arg PORT "${stackPfx}Port" '.[] | select(.OutputKey==$PORT)|.OutputValue' outputs.json )"
address="$(jq -r --arg HOST "${stackPfx}Host" '.[] | select(.OutputKey==$HOST)|.OutputValue' outputs.json )"
admin_username="$(jq -r --arg USER "${stackPfx}Adminusername" '.[] | select(.OutputKey==$USER)|.OutputValue' outputs.json )"
password_arn="$(jq -r --arg ARN "${stackPfx}Adminpasswordarn" '.[] | select(.OutputKey==$ARN)|.OutputValue' outputs.json )"
admin_password="$(aws --output json secretsmanager get-secret-value --secret-id "${password_arn}" --query 'SecretString' | jq -r .|jq -r .password)"

cat > /run/secrets/output<<EOF
services: rds: {
    default: true
    address: "${address}"
    secrets: ["admin"]
    ports: [${port}]
    data: dbname: "${DATABASE_NAME}"
}

secrets: "admin": {
	type: "basic"
	data: {
        username: "${admin_username}"
        password: "${admin_password}"
    }
}
EOF