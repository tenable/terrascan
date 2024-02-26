/*
    Copyright (C) 2022 Tenable, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package config

import (
	"github.com/awslabs/goformation/v7/cloudformation/rds"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// DBInstanceConfig holds config for aws_db_instance
type DBInstanceConfig struct {
	Config
	EnabledCloudWatchLogsExports []string `json:"enabled_cloudwatch_logs_exports"`
	AutoMinorVersionUpgrade      bool     `json:"auto_minor_version_upgrade"`
	CaCertIdentifier             string   `json:"ca_cert_identifier"`
	StorageEncrypted             bool     `json:"storage_encrypted"`
	KmsKeyID                     string   `json:"kms_key_id,omitempty"`
	IamDBAuthEnabled             bool     `json:"iam_database_authentication_enabled"`
	PubliclyAccessible           bool     `json:"publicly_accessible"`
	BackupRetentionPeriod        int      `json:"backup_retention_period"`
	Username                     string   `json:"username"`
	Password                     string   `json:"password"`
	InstanceClass                string   `json:"instance_class"`
	Engine                       string   `json:"engine"`
	EngineVersion                string   `json:"engine_version"`
	Identifier                   string   `json:"identifier"`
	StorageType                  string   `json:"storage_type"`
	DeleteAutomatedBackups       bool     `json:"delete_automated_backups"`
	DeletionProtection           bool     `json:"deletion_protection"`
}

// GetDBInstanceConfig returns config for aws_db_instance
// aws_db_instance
func GetDBInstanceConfig(d *rds.DBInstance) []AWSResourceConfig {
	cf := DBInstanceConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(d.Tags),
			Name: functions.GetVal(d.DBName),
		},
		EnabledCloudWatchLogsExports: d.EnableCloudwatchLogsExports,
		AutoMinorVersionUpgrade:      functions.GetVal(d.AutoMinorVersionUpgrade),
		StorageEncrypted:             functions.GetVal(d.StorageEncrypted),
		KmsKeyID:                     functions.GetVal(d.KmsKeyId),
		CaCertIdentifier:             functions.GetVal(d.CACertificateIdentifier),
		IamDBAuthEnabled:             functions.GetVal(d.EnableIAMDatabaseAuthentication),
		PubliclyAccessible:           functions.GetVal(d.PubliclyAccessible),
		BackupRetentionPeriod:        functions.GetVal(d.BackupRetentionPeriod),
		Username:                     functions.GetVal(d.MasterUsername),
		Password:                     functions.GetVal(d.MasterUserPassword),
		InstanceClass:                functions.GetVal(d.DBInstanceClass),
		Engine:                       functions.GetVal(d.Engine),
		EngineVersion:                functions.GetVal(d.EngineVersion),
		Identifier:                   functions.GetVal(d.DBInstanceIdentifier),
		StorageType:                  functions.GetVal(d.StorageType),
		DeleteAutomatedBackups:       functions.GetVal(d.DeleteAutomatedBackups),
		DeletionProtection:           functions.GetVal(d.DeletionProtection),
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: d.AWSCloudFormationMetadata,
	}}
}
