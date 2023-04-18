#!/bin/sh
set -e

cdk synth --path-metadata false --lookups false > cfn.yaml

mv cfn.yaml /acorn/data/cfn.yaml
cat /acorn/data/cfn.yaml