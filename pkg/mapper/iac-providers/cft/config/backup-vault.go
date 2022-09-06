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
	"github.com/awslabs/goformation/v6/cloudformation/backup"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// BackupVaultConfig holds config for BackupVault
type BackupVaultConfig struct {
	Config
	Name      string `json:"name"`
	KMSKeyARN string `json:"kms_key_arn"`
}

// GetBackupVaultConfig returns config for BackupVault
func GetBackupVaultConfig(b *backup.BackupVault) []AWSResourceConfig {
	cf := BackupVaultConfig{
		Config: Config{
			Name: b.BackupVaultName,
			Tags: b.BackupVaultTags,
		},
		Name:      b.BackupVaultName,
		KMSKeyARN: functions.GetString(b.EncryptionKeyArn),
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: b.AWSCloudFormationMetadata,
	}}
}
