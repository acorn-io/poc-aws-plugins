#!/bin/bash

set -e

address=$(jq -r '.[] | select(.OutputKey=="AMPEndpointURL")|.OutputValue' outputs.json)
arn=$(jq -r '.[]| select(.OutputKey=="AMPWorkspaceArn")|.OutputValue' outputs.json)

cat > /run/secrets/output<<EOF
services: default: {
    default: true
    address: "${address}"
    secrets: ["amp-config"]
}

secrets: "amp-config": {
    type: "opaque"
    data: arn: "${arn}"
}
EOF