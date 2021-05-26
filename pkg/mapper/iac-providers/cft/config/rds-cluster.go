package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/rds"
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
		},
		BackupRetentionPeriod: c.BackupRetentionPeriod,
		StorageEncrypted:      c.StorageEncrypted,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
