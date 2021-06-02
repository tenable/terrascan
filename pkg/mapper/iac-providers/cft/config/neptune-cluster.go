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
	"github.com/awslabs/goformation/v4/cloudformation/neptune"
)

// NeptuneClusterConfig holds config for aws_neptune_cluster
type NeptuneClusterConfig struct {
	Config
	EnableCloudwatchLogsExports []string `json:"enable_cloudwatch_logs_exports,omitempty"`
	StorageEncrypted            bool     `json:"storage_encrypted,omitempty"`
}

// GetNeptuneClusterConfig returns config for aws_neptune_cluster
func GetNeptuneClusterConfig(d *neptune.DBCluster) []AWSResourceConfig {
	cf := NeptuneClusterConfig{
		Config: Config{
			Tags: d.Tags,
		},
		StorageEncrypted:            d.StorageEncrypted,
		EnableCloudwatchLogsExports: d.EnableCloudwatchLogsExports,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
