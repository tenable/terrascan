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

// SecretRotationRulesBlock holds config for SecretRotationRules
type SecretRotationRulesBlock struct {
	AutomaticallyAfterDays int `json:"automatically_after_days"`
}

// SecretsManagerSecretRotationConfig holds config for SecretsManagerSecretRotation
type SecretsManagerSecretRotationConfig struct {
	Config
	SecretID          string                     `json:"secret_id"`
	RotationLambdaARN string                     `json:"rotation_lambda_arn"`
	RotationRules     []SecretRotationRulesBlock `json:"rotation_rules"`
}

// GetSecretsManagerSecretRotationConfig returns config for SecretsManagerSecretRotation
// aws_secretsmanager_secret_rotation no policy
func GetSecretsManagerSecretRotationConfig(r *secretsmanager.RotationSchedule) []AWSResourceConfig {
	var rotationRules []SecretRotationRulesBlock
	if r.RotationRules != nil {
		rotationRules = make([]SecretRotationRulesBlock, 1)
		rotationRules[0].AutomaticallyAfterDays = functions.GetVal(r.RotationRules.AutomaticallyAfterDays)
	}

	cf := SecretsManagerSecretRotationConfig{
		SecretID:          r.SecretId,
		RotationLambdaARN: functions.GetVal(r.RotationLambdaARN),
		RotationRules:     rotationRules,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: r.AWSCloudFormationMetadata,
	}}
}
