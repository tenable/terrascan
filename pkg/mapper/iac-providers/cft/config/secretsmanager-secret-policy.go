package config

import (
	"encoding/json"

	"github.com/awslabs/goformation/v4/cloudformation/secretsmanager"
)

// SecretsManagerSecretPolicyConfig holds config for aws_secretsmanager_secret_policy
type SecretsManagerSecretPolicyConfig struct {
	Config
	ResourcePolicy string `json:"policy"`
}

// GetSecretsManagerSecretPolicyConfig returns config for aws_secretsmanager_secret_policy
func GetSecretsManagerSecretPolicyConfig(s *secretsmanager.ResourcePolicy) []AWSResourceConfig {
	cf := SecretsManagerSecretPolicyConfig{}
	policy, err := json.Marshal(s.ResourcePolicy)
	if err == nil {
		cf.ResourcePolicy = string(policy)
	}
	return []AWSResourceConfig{{Resource: cf}}
}
