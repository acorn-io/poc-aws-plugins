#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { BackupStack } from '../lib/backup-stack';

const app = new cdk.App();

const synthesizer = new cdk.DefaultStackSynthesizer({ generateBootstrapVersionRule: false });

new BackupStack(app, 'BackupStack', {
  env: { account: process.env.CDK_DEFAULT_ACCOUNT, region: process.env.CDK_DEFAULT_REGION },
  synthesizer: synthesizer
});