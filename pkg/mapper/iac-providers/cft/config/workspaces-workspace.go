package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/workspaces"
)

// WorkspacesWorkspaceConfig holds config for aws_workspaces_workspace
type WorkspacesWorkspaceConfig struct {
	Config
	RootVolumeEncryptionEnabled bool `json:"root_volume_encryption_enabled,omitempty"`
	UserVolumeEncryptionEnabled bool `json:"user_volume_encryption_enabled,omitempty"`
}

// GetWorkspacesWorkspaceConfig returns config for aws_workspaces_workspace
func GetWorkspacesWorkspaceConfig(w *workspaces.Workspace) []AWSResourceConfig {
	cf := WorkspacesWorkspaceConfig{
		Config: Config{
			Tags: w.Tags,
		},
		UserVolumeEncryptionEnabled: w.UserVolumeEncryptionEnabled,
		RootVolumeEncryptionEnabled: w.RootVolumeEncryptionEnabled,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
