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
	"github.com/awslabs/goformation/v7/cloudformation/codebuild"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// ArtifactBlock holds config for ArtifactBlock
type ArtifactBlock struct {
	Type string `json:"type"`
}

// CacheBlock holds config for CacheBlock
type CacheBlock struct {
	Type  string   `json:"type"`
	Modes []string `json:"modes"`
}

// CodebuildEnvironmentBlock holds config for CodebuildEnvironmentBlock
type CodebuildEnvironmentBlock struct {
	ComputeType              string `json:"compute_type"`
	Image                    string `json:"image"`
	Type                     string `json:"type"`
	ImagePullCredentialsType string `json:"image_pull_credentials_type"`
}

// SourceBlock holds config for SourceBlock
type SourceBlock struct {
	Type          string `json:"type"`
	Location      string `json:"location"`
	GitCloneDepth int    `json:"git_clone_depth"`
}

// CodebuildProjectConfig holds config for CodebuildProject
type CodebuildProjectConfig struct {
	Config
	Name          string                      `json:"name"`
	Description   string                      `json:"description"`
	BuildTimeout  int                         `json:"build_timeout"`
	QueuedTimeout int                         `json:"queued_timeout"`
	ServiceRole   string                      `json:"service_role"`
	EncryptionKey string                      `json:"encryption_key"`
	Artifacts     []ArtifactBlock             `json:"artifacts"`
	Cache         []CacheBlock                `json:"cache"`
	Environment   []CodebuildEnvironmentBlock `json:"environment"`
	Source        []SourceBlock               `json:"source"`
}

// GetCodebuildProjectConfig returns CodebuildProject
// aws_codebuild_project
func GetCodebuildProjectConfig(p *codebuild.Project) []AWSResourceConfig {

	var artifactBlock []ArtifactBlock
	if p.Artifacts != nil {
		artifactBlock = make([]ArtifactBlock, 1)

		artifactBlock[0].Type = p.Artifacts.Type
	}

	var cacheBlock []CacheBlock
	if p.Cache != nil {
		cacheBlock = make([]CacheBlock, 1)

		cacheBlock[0].Type = p.Cache.Type
		cacheBlock[0].Modes = p.Cache.Modes
	}

	var environmentBlock []CodebuildEnvironmentBlock
	if p.Environment != nil {
		environmentBlock = make([]CodebuildEnvironmentBlock, 1)

		environmentBlock[0].ComputeType = p.Environment.ComputeType
		environmentBlock[0].Image = p.Environment.Image
		environmentBlock[0].Type = p.Environment.Type
		environmentBlock[0].ImagePullCredentialsType = functions.GetVal(p.Environment.ImagePullCredentialsType)
	}

	var sourceBlock []SourceBlock
	if p.Source != nil {
		sourceBlock = make([]SourceBlock, 1)

		sourceBlock[0].Type = p.Source.Type
		sourceBlock[0].Location = functions.GetVal(p.Source.Location)
		sourceBlock[0].GitCloneDepth = functions.GetVal(p.Source.GitCloneDepth)
	}

	cf := CodebuildProjectConfig{
		Config: Config{
			Name: functions.GetVal(p.Name),
		},
		Name:          functions.GetVal(p.Name),
		Description:   functions.GetVal(p.Description),
		BuildTimeout:  functions.GetVal(p.TimeoutInMinutes),
		QueuedTimeout: functions.GetVal(p.QueuedTimeoutInMinutes),
		ServiceRole:   p.ServiceRole,
		EncryptionKey: functions.GetVal(p.EncryptionKey),
		Artifacts:     artifactBlock,
		Cache:         cacheBlock,
		Environment:   environmentBlock,
		Source:        sourceBlock,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: p.AWSCloudFormationMetadata,
	}}
}
