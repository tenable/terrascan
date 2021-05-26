package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/secretsmanager"
)

// SecretsManagerSecretConfig holds config for aws_secretsmanager_secret
type SecretsManagerSecretConfig struct {
	Config
	KmsKeyID string `json:"kms_key_id,omitempty"`
}

// GetSecretsManagerSecretConfig returns config for aws_secretsmanager_secret
func GetSecretsManagerSecretConfig(s *secretsmanager.Secret) []AWSResourceConfig {
	cf := SecretsManagerSecretConfig{
		Config: Config{
			Tags: s.Tags,
		},
		KmsKeyID: s.KmsKeyId,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
