package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/logs"
)

// LogCloudWatchGroupConfig holds config for aws_cloudwatch_log_group
type LogCloudWatchGroupConfig struct {
	Config
	LogGroupName    string `json:"name"`
	KmsKeyID        string `json:"kms_key_id,omitempty"`
	RetentionInDays int    `json:"retention_in_days"`
}

// GetLogCloudWatchGroupConfig returns config for aws_cloudwatch_log_group
func GetLogCloudWatchGroupConfig(r *logs.LogGroup) []AWSResourceConfig {
	cf := LogCloudWatchGroupConfig{
		Config: Config{
			Name: r.LogGroupName,
		},
		LogGroupName:    r.LogGroupName,
		KmsKeyID:        r.KmsKeyId,
		RetentionInDays: r.RetentionInDays,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
