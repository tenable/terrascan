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
	"github.com/awslabs/goformation/v7/cloudformation/secretsmanager"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// SecretsManagerSecretConfig holds config for aws_secretsmanager_secret
type SecretsManagerSecretConfig struct {
	Config
	KmsKeyID string `json:"kms_key_id,omitempty"`
}

// GetSecretsManagerSecretConfig returns config for aws_secretsmanager_secret
// aws_secretsmanager_secret
func GetSecretsManagerSecretConfig(s *secretsmanager.Secret) []AWSResourceConfig {
	cf := SecretsManagerSecretConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(s.Tags),
			Name: functions.GetVal(s.Name),
		},
		KmsKeyID: functions.GetVal(s.KmsKeyId),
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: s.AWSCloudFormationMetadata,
	}}
}
