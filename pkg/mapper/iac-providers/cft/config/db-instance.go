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
	"github.com/awslabs/goformation/v6/cloudformation/rds"
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
}

// GetDBInstanceConfig returns config for aws_db_instance
func GetDBInstanceConfig(d *rds.DBInstance) []AWSResourceConfig {
	cf := DBInstanceConfig{
		Config: Config{
			Tags: d.Tags,
			Name: functions.GetVal(d.DBName),
		},
		EnabledCloudWatchLogsExports: functions.GetVal(d.EnableCloudwatchLogsExports),
		AutoMinorVersionUpgrade:      functions.GetVal(d.AutoMinorVersionUpgrade),
		StorageEncrypted:             functions.GetVal(d.StorageEncrypted),
		KmsKeyID:                     functions.GetVal(d.KmsKeyId),
		CaCertIdentifier:             functions.GetVal(d.CACertificateIdentifier),
		IamDBAuthEnabled:             functions.GetVal(d.EnableIAMDatabaseAuthentication),
		PubliclyAccessible:           functions.GetVal(d.PubliclyAccessible),
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: d.AWSCloudFormationMetadata,
	}}
}
