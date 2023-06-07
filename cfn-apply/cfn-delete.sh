#!/bin/bash

set -e

if [ $# -eq 0 ]; then
  echo "No arguments passed..."
  exit 1
fi 

stack="${1}"
sanitizedName=$(sed 's/-//g' <<< ${stack})

if [ -z "${ACORN_EVENT}" ]; then
  echo "Acorn event not set..."
  exit 0
fi

if [ "${ACORN_EVENT}" = "delete" ]; then
    aws cloudformation delete-stack --stack-name "${sanitizedName}"
    exit 0
fi

echo "No event found to process..."
