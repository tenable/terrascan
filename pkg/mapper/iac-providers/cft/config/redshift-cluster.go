package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/redshift"
)

// RedshiftClusterConfig holds config for aws_redshift_cluster
type RedshiftClusterConfig struct {
	Config
	LoggingProperties  interface{} `json:"logging,omitempty"`
	KmsKeyID           string      `json:"kms_key_id,omitempty"`
	Encrypted          bool        `json:"encrypted"`
	PubliclyAccessible bool        `json:"publicly_accessible"`
}

// GetRedshiftClusterConfig returns config for aws_redshift_cluster
func GetRedshiftClusterConfig(c *redshift.Cluster) []AWSResourceConfig {
	cf := RedshiftClusterConfig{
		Config: Config{
			Name: c.DBName,
		},
		KmsKeyID:           c.KmsKeyId,
		Encrypted:          c.Encrypted,
		PubliclyAccessible: c.PubliclyAccessible,
	}
	if c.LoggingProperties != nil {
		// if LoggingProperties are mentioned in cft,
		// its always enabled
		logging := make(map[string]bool)
		logging["enable"] = true
		cf.LoggingProperties = []map[string]bool{logging}
	}
	return []AWSResourceConfig{{Resource: cf}}
}
