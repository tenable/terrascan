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
	"github.com/awslabs/goformation/v7/cloudformation/workspaces"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// WorkspacesWorkspaceConfig holds config for aws_workspaces_workspace
type WorkspacesWorkspaceConfig struct {
	Config
	RootVolumeEncryptionEnabled bool `json:"root_volume_encryption_enabled,omitempty"`
	UserVolumeEncryptionEnabled bool `json:"user_volume_encryption_enabled,omitempty"`
}

// GetWorkspacesWorkspaceConfig returns config for aws_workspaces_workspace
// aws_workspaces_workspace
func GetWorkspacesWorkspaceConfig(w *workspaces.Workspace) []AWSResourceConfig {
	cf := WorkspacesWorkspaceConfig{
		Config: Config{
			Name: w.UserName,
			Tags: w.Tags,
		},
		UserVolumeEncryptionEnabled: functions.GetVal(w.UserVolumeEncryptionEnabled),
		RootVolumeEncryptionEnabled: functions.GetVal(w.RootVolumeEncryptionEnabled),
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: w.AWSCloudFormationMetadata,
	}}
}
