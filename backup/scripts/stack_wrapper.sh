#!/bin/bash

set -e

cfn_script="/app/cfn-apply.sh"
if [ "${ACORN_EVENT}" = "delete" ]; then
  cfn_script="/app/cfn-delete.sh"
fi

APP=$(sed 's/-//g' <<< ${ACORN_APP^})
NS=$(sed 's/-//g' <<< ${ACORN_NAMESPACE^})

stackName="${APP}${NS}${DB}Stack"

exec ${cfn_script} ${stackName}