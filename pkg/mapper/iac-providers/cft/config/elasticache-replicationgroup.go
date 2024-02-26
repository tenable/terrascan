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
	"github.com/awslabs/goformation/v7/cloudformation/elasticache"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// ElastiCacheReplicationGroupConfig holds config for aws_elasticache_replication_group
type ElastiCacheReplicationGroupConfig struct {
	Config
	AtRestEncryptionEnabled  bool `json:"at_rest_encryption_enabled,omitempty"`
	TransitEncryptionEnabled bool `json:"transit_encryption_enabled,omitempty"`
}

// GetElastiCacheReplicationGroupConfig returns config for aws_elasticache_replication_group
// aws_elasticache_replication_group
func GetElastiCacheReplicationGroupConfig(r *elasticache.ReplicationGroup) []AWSResourceConfig {
	cf := ElastiCacheReplicationGroupConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(r.Tags),
		},
		AtRestEncryptionEnabled:  functions.GetVal(r.AtRestEncryptionEnabled),
		TransitEncryptionEnabled: functions.GetVal(r.TransitEncryptionEnabled),
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: r.AWSCloudFormationMetadata,
	}}
}
