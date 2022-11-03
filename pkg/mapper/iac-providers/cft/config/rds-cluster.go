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
			Name: functions.GetVal(c.DatabaseName),
			Tags: c.Tags,
		},
		BackupRetentionPeriod: functions.GetVal(c.BackupRetentionPeriod),
		StorageEncrypted:      functions.GetVal(c.StorageEncrypted),
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: c.AWSCloudFormationMetadata,
	}}
}
