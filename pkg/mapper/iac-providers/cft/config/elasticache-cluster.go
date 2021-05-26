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
