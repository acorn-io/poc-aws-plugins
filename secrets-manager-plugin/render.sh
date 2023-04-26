#!/bin/bash

set -e

arn="${1}"

data="$(aws --output json secretsmanager get-secret-value --secret-id "${arn}" --query 'SecretString' | jq -r .|jq '.'|sed 's/,$//')"

cat > /run/secrets/output<<EOF
services: "secret-manager": {
    secrets: ["secret-value"]
}

secrets: "secret-value": {
    type: "opaque"
    data: ${data}
}
EOF