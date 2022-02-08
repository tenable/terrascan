/*
    Copyright (C) 2022 Accurics, Inc.

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

import "github.com/awslabs/goformation/v5/cloudformation/redshift"

// ParameterBlock holds config for Parameter
type ParameterBlock struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// RedshiftParameterGroupConfig holds config for RedshiftParameterGroup
type RedshiftParameterGroupConfig struct {
	Config
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Family      string           `json:"family"`
	Parameter   []ParameterBlock `json:"parameter"`
}

// GetRedshiftParameterGroupConfig returns config for RedshiftParameterGroup
func GetRedshiftParameterGroupConfig(p *redshift.ClusterParameterGroup, paramGroupName string) []AWSResourceConfig {
	parameterBlock := make([]ParameterBlock, len(p.Parameters))
	for i := range p.Parameters {
		parameterBlock[i].Name = p.Parameters[i].ParameterName
		parameterBlock[i].Value = p.Parameters[i].ParameterValue
	}

	cf := RedshiftParameterGroupConfig{
		Config: Config{
			Name: paramGroupName,
			Tags: p.Tags,
		},
		Name:        paramGroupName,
		Description: p.Description,
		Family:      p.ParameterGroupFamily,
		Parameter:   parameterBlock,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: p.AWSCloudFormationMetadata,
	}}
}
