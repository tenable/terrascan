package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/elasticache"
)

// ElastiCacheReplicationGroupConfig holds config for aws_elasticache_replication_group
type ElastiCacheReplicationGroupConfig struct {
	Config
	AtRestEncryptionEnabled  bool `json:"at_rest_encryption_enabled,omitempty"`
	TransitEncryptionEnabled bool `json:"transit_encryption_enabled,omitempty"`
}

// GetElastiCacheReplicationGroupConfig returns config for aws_elasticache_replication_group
func GetElastiCacheReplicationGroupConfig(r *elasticache.ReplicationGroup) []AWSResourceConfig {
	cf := ElastiCacheReplicationGroupConfig{
		Config: Config{
			Tags: r.Tags,
		},
		AtRestEncryptionEnabled:  r.AtRestEncryptionEnabled,
		TransitEncryptionEnabled: r.TransitEncryptionEnabled,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
