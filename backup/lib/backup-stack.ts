import * as cdk from 'aws-cdk-lib';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as backup from 'aws-cdk-lib/aws-backup';
import { Construct } from 'constructs';
import { Key } from 'aws-cdk-lib/aws-kms';

export class BackupStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const uid = cdk.Fn.select(0, cdk.Fn.split('-', cdk.Fn.select(2, cdk.Fn.split('/', this.stackId))))

    const backupRole = new iam.Role(this, 'Role', {
      assumedBy: new iam.ServicePrincipal('backup.amazonaws.com'),
      description: 'Role allows AWS Backup service required access to other services',
      managedPolicies: [iam.ManagedPolicy.fromAwsManagedPolicyName("service-role/AWSBackupServiceRolePolicyForBackup")]
    });

    const backupVault = new backup.BackupVault(this, 'BackupVault', {
      backupVaultName: `acorn_backup_${uid}`,
      removalPolicy: cdk.RemovalPolicy.DESTROY
    })

    const plan = backup.BackupPlan.daily35DayRetention(this, 'Plan', backupVault);

    const backupSelection = new backup.BackupSelection(this, 'BackupSelection', {
      backupPlan: plan,
      allowRestores: true,
      backupSelectionName: `acorn_backups_${uid}`,
      resources: [
        backup.BackupResource.fromTag("acorn.io/backup", "true", backup.TagOperation.STRING_EQUALS)

      ],
      role: backupRole,
    });

  }
}