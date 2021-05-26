package config

import (
	"encoding/json"

	"github.com/awslabs/goformation/v4/cloudformation/efs"
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
func GetEfsFileSystemConfig(f *efs.FileSystem) []AWSResourceConfig {

	resourceConfigs := make([]AWSResourceConfig, 0)

	resourceConfigs = append(resourceConfigs, AWSResourceConfig{
		Resource: EfsFileSystemConfig{
			Config:    Config{},
			KmsKeyID:  f.KmsKeyId,
			Encrypted: f.Encrypted,
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
			Type:     EfsFileSystemPolicy,
			Name:     "efs",
		})
	}

	return resourceConfigs
}
