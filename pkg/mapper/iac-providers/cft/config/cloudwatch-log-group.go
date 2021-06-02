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
	"github.com/awslabs/goformation/v4/cloudformation/logs"
)

// LogCloudWatchGroupConfig holds config for aws_cloudwatch_log_group
type LogCloudWatchGroupConfig struct {
	Config
	LogGroupName    string `json:"name"`
	KmsKeyID        string `json:"kms_key_id,omitempty"`
	RetentionInDays int    `json:"retention_in_days"`
}

// GetLogCloudWatchGroupConfig returns config for aws_cloudwatch_log_group
func GetLogCloudWatchGroupConfig(r *logs.LogGroup) []AWSResourceConfig {
	cf := LogCloudWatchGroupConfig{
		Config: Config{
			Name: r.LogGroupName,
		},
		LogGroupName:    r.LogGroupName,
		KmsKeyID:        r.KmsKeyId,
		RetentionInDays: r.RetentionInDays,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: r.AWSCloudFormationMetadata,
	}}
}
