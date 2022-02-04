/*
    Copyright (C) 2021 Accurics, Inc.

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
	"github.com/awslabs/goformation/v5/cloudformation/rds"
)

// RDSClusterConfig holds config for aws_rds_cluster
type RDSClusterConfig struct {
	Config
	BackupRetentionPeriod int  `json:"backup_retention_period,omitempty"`
	StorageEncrypted      bool `json:"storage_encrypted"`
}

// GetRDSClusterConfig returns config for aws_rds_cluster
func GetRDSClusterConfig(c *rds.DBCluster) []AWSResourceConfig {
	cf := RDSClusterConfig{
		Config: Config{
			Name: c.DatabaseName,
			Tags: c.Tags,
		},
		BackupRetentionPeriod: c.BackupRetentionPeriod,
		StorageEncrypted:      c.StorageEncrypted,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: c.AWSCloudFormationMetadata,
	}}
}
