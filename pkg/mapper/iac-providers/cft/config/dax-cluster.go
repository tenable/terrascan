package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/dax"
)

// DaxClusterConfig holds config for aws_dax_cluster
type DaxClusterConfig struct {
	Config
	ServerSideEncryption []map[string]interface{} `json:"server_side_encryption"`
}

// GetDaxClusterConfig returns config for aws_dax_cluster
func GetDaxClusterConfig(t *dax.Cluster) []AWSResourceConfig {
	cf := DaxClusterConfig{
		Config: Config{
			Tags: t.Tags,
			Name: t.ClusterName,
		},
	}
	sse := make(map[string]interface{})
	if t.SSESpecification != nil {
		sse["enabled"] = t.SSESpecification.SSEEnabled
	}
	cf.ServerSideEncryption = []map[string]interface{}{sse}
	return []AWSResourceConfig{{Resource: cf}}
}
