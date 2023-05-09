#!/bin/bash

set -e

cat > /run/secrets/output<<EOF
services: default: {
    default: true
}
EOF
