/*
    Copyright (C) 2022 Tenable, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package config

import (
	"encoding/json"

	"github.com/awslabs/goformation/v7/cloudformation/secretsmanager"
)

// SecretsManagerSecretPolicyConfig holds config for aws_secretsmanager_secret_policy
type SecretsManagerSecretPolicyConfig struct {
	Config
	ResourcePolicy string `json:"policy"`
}

// GetSecretsManagerSecretPolicyConfig returns config for aws_secretsmanager_secret_policy
// aws_secretsmanager_secret_policy
func GetSecretsManagerSecretPolicyConfig(s *secretsmanager.ResourcePolicy) []AWSResourceConfig {
	cf := SecretsManagerSecretPolicyConfig{}
	policy, err := json.Marshal(s.ResourcePolicy)
	if err == nil {
		cf.ResourcePolicy = string(policy)
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: s.AWSCloudFormationMetadata,
	}}
}
