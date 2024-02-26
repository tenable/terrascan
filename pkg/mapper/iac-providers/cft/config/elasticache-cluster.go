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

// ElastiCacheClusterConfig holds config for aws_elasticache_cluster
type ElastiCacheClusterConfig struct {
	Config
	AZMode        string `json:"az_mode"`
	Engine        string `json:"engine,omitempty"`
	EngineVersion string `json:"engine_version,omitempty"`
}

// GetElastiCacheClusterConfig returns config for aws_elasticache_cluster
// aws_elasticache_cluster
func GetElastiCacheClusterConfig(e *elasticache.CacheCluster) []AWSResourceConfig {
	cf := ElastiCacheClusterConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(e.Tags),
			Name: functions.GetVal(e.ClusterName),
		},
		AZMode:        functions.GetVal(e.AZMode),
		Engine:        e.Engine,
		EngineVersion: functions.GetVal(e.EngineVersion),
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: e.AWSCloudFormationMetadata,
	}}
}
