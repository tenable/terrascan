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
	"github.com/awslabs/goformation/v4/cloudformation/config"
)

// ConfigurationAggregatorConfig holds config for aws_config_configuration_aggregator
type ConfigurationAggregatorConfig struct {
	Config
	AccountAggregationSource interface{} `json:"account_aggregation_source,omitempty"`
	OrgAggregationSource     interface{} `json:"organization_aggregation_source,omitempty"`
}

// GetConfigConfigurationAggregatorConfig returns config for aws_config_configuration_aggregator
func GetConfigConfigurationAggregatorConfig(c *config.ConfigurationAggregator) []AWSResourceConfig {
	cf := ConfigurationAggregatorConfig{
		Config: Config{
			Tags: c.Tags,
			Name: c.ConfigurationAggregatorName,
		},
	}
	if c.AccountAggregationSources != nil {
		accountAggregationSources := make([]map[string]interface{}, 0)
		for i := range c.AccountAggregationSources {
			accountAggregationSource := make(map[string]interface{})
			accountAggregationSource["all_regions"] = c.AccountAggregationSources[i].AllAwsRegions
			accountAggregationSources = append(accountAggregationSources, accountAggregationSource)
		}
		if len(accountAggregationSources) > 0 {
			cf.AccountAggregationSource = accountAggregationSources
		}
	}
	if c.OrganizationAggregationSource != nil {
		organizationAggregationSources := make([]map[string]interface{}, 0)
		organizationAggregationSource := make(map[string]interface{})
		organizationAggregationSource["all_regions"] = c.OrganizationAggregationSource.AllAwsRegions
		organizationAggregationSources = append(organizationAggregationSources, organizationAggregationSource)
		if len(organizationAggregationSources) > 0 {
			cf.OrgAggregationSource = organizationAggregationSources
		}
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: c.AWSCloudFormationMetadata,
	}}
}
