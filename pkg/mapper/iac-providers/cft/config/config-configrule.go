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
	"github.com/awslabs/goformation/v7/cloudformation/config"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// AWSConfigConfigRuleConfig holds config for aws_config_config_rule
type AWSConfigConfigRuleConfig struct {
	Config
	Source interface{} `json:"source"`
}

// GetConfigConfigRuleConfig returns config for aws_config_config_rule
// aws_config_config_rule
func GetConfigConfigRuleConfig(c *config.ConfigRule) []AWSResourceConfig {
	cf := AWSConfigConfigRuleConfig{
		Config: Config{
			Name: functions.GetVal(c.ConfigRuleName),
		},
	}
	if c.Source != nil {
		sources := make([]map[string]interface{}, 0)
		source := make(map[string]interface{})
		source["source_identifier"] = functions.GetVal(c.Source.SourceIdentifier)
		sources = append(sources, source)
		if len(sources) > 0 {
			cf.Source = sources
		}
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: c.AWSCloudFormationMetadata,
	}}
}
