#!/bin/bash
set -e -x

APP=$(sed 's/-//g' <<< ${ACORN_APP^})
NS=$(sed 's/-//g' <<< ${ACORN_NAMESPACE^})
DB=$(sed 's/-//g' <<< ${DATABASE_NAME^})

STACK_PREFIX="${APP}${NS}${DB}"
STACK_NAME="${STACK_PREFIX}Stack"

if [ "${ACORN_EVENT}" = "delete" ]; then
    aws cloudformation delete-stack --stack-name "${STACK_NAME}"
    exit 0
fi

# Run CDK synth
cat cdk.context.json
cdk synth --path-metadata false --lookups false > cfn.yaml
cat cfn.yaml

# Run CloudFormation
./scripts/stacklog.sh ${STACK_NAME} &
aws cloudformation deploy --template-file cfn.yaml --stack-name "${STACK_NAME}" --capabilities CAPABILITY_IAM --capabilities CAPABILITY_NAMED_IAM --no-fail-on-empty-changeset --no-cli-pager
aws cloudformation describe-stacks --stack-name "${STACK_NAME}" --query 'Stacks[0].Outputs' > outputs.json

# Render Output
PORT="$(          jq -r '.[] | select(.OutputKey=="'${STACK_PREFIX}'Port")            |.OutputValue' outputs.json )"
ADDRESS="$(       jq -r '.[] | select(.OutputKey=="'${STACK_PREFIX}'Host")            |.OutputValue' outputs.json )"
ADMIN_USERNAME="$(jq -r '.[] | select(.OutputKey=="'${STACK_PREFIX}'Adminusername")   |.OutputValue' outputs.json )"
PASSWORD_ARN="$(  jq -r '.[] | select(.OutputKey=="'${STACK_PREFIX}'Adminpasswordarn")|.OutputValue' outputs.json )"

# Turn off echo
set +x
ADMIN_PASSWORD="$(aws --output json secretsmanager get-secret-value --secret-id "${PASSWORD_ARN}" --query 'SecretString' | jq -r .|jq -r .password)"

cat > /run/secrets/output <<EOF
services: rds: {
    default: true
    address: "${ADDRESS}"
    secrets: ["admin"]
    ports: [${PORT}]
    data: dbName: "${DATABASE_NAME}"
}

secrets: "admin": {
	type: "basic"
	data: {
        username: "${ADMIN_USERNAME}"
        password: "${ADMIN_PASSWORD}"
    }
}
EOF
