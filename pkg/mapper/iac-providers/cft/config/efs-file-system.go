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

	"github.com/awslabs/goformation/v7/cloudformation/efs"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

const (
	// EfsFileSystemPolicy represents the sub-resource aws_efs_file_system_policy for attribute FileSystemPolicy
	EfsFileSystemPolicy = "FileSystemPolicy"
)

// EfsFileSystemConfig holds config for aws_efs_file_system
type EfsFileSystemConfig struct {
	Config
	Encrypted bool   `json:"encrypted"`
	KmsKeyID  string `json:"kms_key_id,omitempty"`
}

// EfsFileSystemPolicyConfig holds config for aws_efs_file_system_policy
type EfsFileSystemPolicyConfig struct {
	Config
	FileSystemPolicy string `json:"policy"`
}

// GetEfsFileSystemConfig returns config for aws_efs_file_system and aws_efs_file_system_policy
// aws_efs_file_system and aws_efs_file_system_policy
func GetEfsFileSystemConfig(f *efs.FileSystem) []AWSResourceConfig {

	resourceConfigs := make([]AWSResourceConfig, 0)

	resourceConfigs = append(resourceConfigs, AWSResourceConfig{
		Metadata: f.AWSCloudFormationMetadata,
		Resource: EfsFileSystemConfig{
			KmsKeyID:  functions.GetVal(f.KmsKeyId),
			Encrypted: functions.GetVal(f.Encrypted),
		},
	})

	if f.FileSystemPolicy != nil {
		policyConfig := EfsFileSystemPolicyConfig{}
		policies, err := json.Marshal(f.FileSystemPolicy)
		if err == nil {
			policyConfig.FileSystemPolicy = string(policies)
		}
		resourceConfigs = append(resourceConfigs, AWSResourceConfig{
			Resource: policyConfig,
			Metadata: f.AWSCloudFormationMetadata,
			Type:     EfsFileSystemPolicy,
			Name:     "efs",
		})
	}

	return resourceConfigs
}
