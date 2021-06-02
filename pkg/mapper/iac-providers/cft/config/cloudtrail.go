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
	"github.com/awslabs/goformation/v4/cloudformation/cloudtrail"
)

// CloudTrailConfig holds config for aws_cloudtrail
type CloudTrailConfig struct {
	Config
	IsMultiRegionTrail      interface{} `json:"is_multi_region_trail"`
	KmsKeyID                interface{} `json:"kms_key_id"`
	SnsTopicName            interface{} `json:"sns_topic_name"`
	EnableLogFileValidation interface{} `json:"enable_log_file_validation"`
}

// GetCloudTrailConfig returns config for aws_cloudtrail
func GetCloudTrailConfig(t *cloudtrail.Trail) []AWSResourceConfig {
	cf := CloudTrailConfig{
		Config:                  Config{Tags: t.Tags, Name: t.TrailName},
		EnableLogFileValidation: t.EnableLogFileValidation,
		IsMultiRegionTrail:      t.IsMultiRegionTrail,
	}
	if len(t.KMSKeyId) > 0 {
		cf.KmsKeyID = t.KMSKeyId
	}
	if len(t.SnsTopicName) > 0 {
		cf.SnsTopicName = t.SnsTopicName
	}

	return []AWSResourceConfig{{Resource: cf}}
}
