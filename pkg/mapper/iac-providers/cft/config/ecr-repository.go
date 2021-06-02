/*
    Copyright (C) 2021 Accurics, Inc.

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
	"github.com/awslabs/goformation/v4/cloudformation/ecr"
)

// EcrRepositoryConfig holds config for aws_ecr_repository
type EcrRepositoryConfig struct {
	Config
	ISC  []map[string]interface{} `json:"image_scanning_configuration"`
	AERP interface{}              `json:"aws_ecr_repository_policy,omitempty"`
}

// ISC holds config for image_scanning_configuration attribute
type ISC struct {
	ScanOnPush bool `json:"scan_on_push"`
}

// GetEcrRepositoryConfig returns config for aws_ecr_repository
func GetEcrRepositoryConfig(r *ecr.Repository) []AWSResourceConfig {
	cf := EcrRepositoryConfig{
		Config: Config{
			Tags: r.Tags,
			Name: r.RepositoryName,
		},
	}
	if r.ImageScanningConfiguration != nil {
		m := r.ImageScanningConfiguration.(map[string]interface{})
		if m["ScanOnPush"] != nil {
			sop := make(map[string]interface{})
			sop["scan_on_push"] = m["ScanOnPush"]
			cf.ISC = []map[string]interface{}{sop}
		} else {
			cf.ISC = make([]map[string]interface{}, 0)
		}
	} else {
		cf.ISC = make([]map[string]interface{}, 0)
	}
	cf.AERP = r.RepositoryPolicyText
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: r.AWSCloudFormationMetadata,
	}}
}
