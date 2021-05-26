package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/neptune"
)

// NeptuneClusterConfig holds config for aws_neptune_cluster
type NeptuneClusterConfig struct {
	Config
	EnableCloudwatchLogsExports []string `json:"enable_cloudwatch_logs_exports,omitempty"`
	StorageEncrypted            bool     `json:"storage_encrypted,omitempty"`
}

// GetNeptuneClusterConfig returns config for aws_neptune_cluster
func GetNeptuneClusterConfig(d *neptune.DBCluster) []AWSResourceConfig {
	cf := NeptuneClusterConfig{
		Config: Config{
			Tags: d.Tags,
		},
		StorageEncrypted:            d.StorageEncrypted,
		EnableCloudwatchLogsExports: d.EnableCloudwatchLogsExports,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
