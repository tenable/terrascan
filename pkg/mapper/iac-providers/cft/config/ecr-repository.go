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
	"github.com/awslabs/goformation/v7/cloudformation/ecr"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// EcrRepositoryConfig holds config for aws_ecr_repository
type EcrRepositoryConfig struct {
	Config
	ImageScanningConfiguration []ImageScanningConfigurationBlock `json:"image_scanning_configuration"`
	AERP                       interface{}                       `json:"aws_ecr_repository_policy,omitempty"`
}

// ImageScanningConfigurationBlock holds config for image_scanning_configuration attribute
type ImageScanningConfigurationBlock struct {
	ScanOnPush bool `json:"scan_on_push"`
}

// GetEcrRepositoryConfig returns config for aws_ecr_repository
// aws_ecr_repository
func GetEcrRepositoryConfig(r *ecr.Repository) []AWSResourceConfig {
	var imageScanningConfiguration []ImageScanningConfigurationBlock
	if r.ImageScanningConfiguration != nil {
		imageScanningConfiguration = make([]ImageScanningConfigurationBlock, 1)
		imageScanningConfiguration[0].ScanOnPush = functions.GetVal(r.ImageScanningConfiguration.ScanOnPush)
	}

	cf := EcrRepositoryConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(r.Tags),
			Name: functions.GetVal(r.RepositoryName),
		},
		ImageScanningConfiguration: imageScanningConfiguration,
	}

	cf.AERP = r.RepositoryPolicyText
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: r.AWSCloudFormationMetadata,
	}}
}
