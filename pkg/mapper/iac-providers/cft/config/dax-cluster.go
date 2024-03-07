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
	"github.com/awslabs/goformation/v7/cloudformation/dax"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// DaxClusterConfig holds config for aws_dax_cluster
type DaxClusterConfig struct {
	Config
	ServerSideEncryption []SSE `json:"server_side_encryption"`
}

// GetDaxClusterConfig returns config for aws_dax_cluster
// aws_dax_cluster
func GetDaxClusterConfig(t *dax.Cluster) []AWSResourceConfig {
	cf := DaxClusterConfig{
		Config: Config{
			Tags: t.Tags,
			Name: functions.GetVal(t.ClusterName),
		},
	}

	if t.SSESpecification != nil {
		cf.ServerSideEncryption = make([]SSE, 1)

		cf.ServerSideEncryption[0].Enabled = functions.GetVal(t.SSESpecification.SSEEnabled)
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: t.AWSCloudFormationMetadata,
	}}
}
