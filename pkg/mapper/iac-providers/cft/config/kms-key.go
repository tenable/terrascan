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

	"github.com/awslabs/goformation/v7/cloudformation/kms"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// KmsKeyConfig holds config for aws_kms_key
type KmsKeyConfig struct {
	Config
	Description         string `json:"description"`
	KeyPolicy           string `json:"policy"`
	PendingWindowInDays int    `json:"deletion_window_in_days"`
	Enabled             bool   `json:"is_enabled"`
	EnableKeyRotation   bool   `json:"enable_key_rotation"`
}

// GetKmsKeyConfig returns config for aws_kms_key
// aws_kms_key
func GetKmsKeyConfig(k *kms.Key) []AWSResourceConfig {
	cf := KmsKeyConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(k.Tags),
		},
		Enabled:             functions.GetVal(k.Enabled),
		EnableKeyRotation:   functions.GetVal(k.EnableKeyRotation),
		PendingWindowInDays: functions.GetVal(k.PendingWindowInDays),
	}

	keyPolicy, err := json.Marshal(k.KeyPolicy)
	if err == nil {
		cf.KeyPolicy = string(keyPolicy)
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: k.AWSCloudFormationMetadata,
	}}
}
