#!/bin/bash
set -e -o pipefail

arn="${1}"

echo Getting secret value for "$arn"
data="$(aws --output json secretsmanager get-secret-value --secret-id "${arn}" | jq '{value: .SecretString}')"

cat > /run/secrets/output <<EOF
services: "secret-manager": {
    default: true
    secrets: ["secret-value"]
}

secrets: "item": {
    type: "opaque"
    data: ${data}
}
EOF
