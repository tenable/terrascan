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
	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// EbsVolumeConfig holds config for aws_ebs_volume
type EbsVolumeConfig struct {
	Config
	Encrypted bool   `json:"encrypted,omitempty"`
	KmsKeyID  string `json:"kms_key_id,omitempty"`
}

// GetEbsVolumeConfig returns config for aws_ebs_volume
// aws_ebs_volume
func GetEbsVolumeConfig(v *ec2.Volume) []AWSResourceConfig {
	cf := EbsVolumeConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(v.Tags),
		},
		Encrypted: functions.GetVal(v.Encrypted),
		KmsKeyID:  functions.GetVal(v.KmsKeyId),
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: v.AWSCloudFormationMetadata,
	}}
}
