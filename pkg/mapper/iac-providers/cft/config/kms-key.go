package config

import (
	"encoding/json"

	"github.com/awslabs/goformation/v4/cloudformation/kms"
)

// KmsKeyConfig holds the config for aws_kms_key
type KmsKeyConfig struct {
	Config
	Description         string `json:"description"`
	KeyPolicy           string `json:"policy"`
	PendingWindowInDays int    `json:"deletion_window_in_days"`
	Enabled             bool   `json:"is_enabled"`
	EnableKeyRotation   bool   `json:"enable_key_rotation"`
}

// GetKmsKeyConfig returns config for aws_kms_key
func GetKmsKeyConfig(k *kms.Key) []AWSResourceConfig {
	cf := KmsKeyConfig{
		Config: Config{
			Tags: k.Tags,
		},
		Enabled:             k.Enabled,
		EnableKeyRotation:   k.EnableKeyRotation,
		PendingWindowInDays: k.PendingWindowInDays,
	}

	keyPolicy, err := json.Marshal(k.KeyPolicy)
	if err == nil {
		cf.KeyPolicy = string(keyPolicy)
	}

	return []AWSResourceConfig{{Resource: cf}}
}
