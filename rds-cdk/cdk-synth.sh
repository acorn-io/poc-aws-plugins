#!/bin/sh

find
cdk --app './rds' synth --path-metadata false --lookups false --no-verion-reporting > cfn.yaml

mv cfn.yaml /acorn/data/cfn.yaml
cat /acorn/data/cfn.yaml
