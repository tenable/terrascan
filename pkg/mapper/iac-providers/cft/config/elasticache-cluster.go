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
	"github.com/awslabs/goformation/v4/cloudformation/elasticache"
)

// ElastiCacheClusterConfig holds config for aws_elasticache_cluster
type ElastiCacheClusterConfig struct {
	Config
	AZMode        string `json:"az_mode"`
	Engine        string `json:"engine,omitempty"`
	EngineVersion string `json:"engine_version,omitempty"`
}

// GetElastiCacheClusterConfig returns config for aws_elasticache_cluster
func GetElastiCacheClusterConfig(e *elasticache.CacheCluster) []AWSResourceConfig {
	cf := ElastiCacheClusterConfig{
		Config: Config{
			Tags: e.Tags,
			Name: e.ClusterName,
		},
		AZMode:        e.AZMode,
		Engine:        e.Engine,
		EngineVersion: e.EngineVersion,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
