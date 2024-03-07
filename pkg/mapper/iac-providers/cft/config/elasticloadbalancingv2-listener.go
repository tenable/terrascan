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
	"github.com/awslabs/goformation/v7/cloudformation/elasticloadbalancingv2"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// ElasticLoadBalancingV2ListenerConfig holds config for aws_lb_listener
type ElasticLoadBalancingV2ListenerConfig struct {
	Config
	Protocol      string                `json:"protocol"`
	Port          int                   `json:"port"`
	DefaultAction []DefaultActionConfig `json:"default_action"`
}

// DefaultActionConfig holds config for default_action attribute of aws_lb_listener
type DefaultActionConfig struct {
	RedirectConfig []RedirectConfig `json:"redirect"`
}

// RedirectConfig holds config for redirect attribute of default_action
type RedirectConfig struct {
	Protocol string `json:"protocol"`
	Port     string `json:"port"`
}

// GetElasticLoadBalancingV2ListenerConfig returns config for aws_lb_listener
// aws_lb_listener
func GetElasticLoadBalancingV2ListenerConfig(l *elasticloadbalancingv2.Listener) []AWSResourceConfig {
	// create a listener subresource per DefaultAction defined in cft
	// as only one default action per listener is possible in terraform
	resourceConfigs := make([]AWSResourceConfig, 0)

	for _, action := range l.DefaultActions {
		// DefaultActions are required
		cf := ElasticLoadBalancingV2ListenerConfig{
			Config:   Config{},
			Protocol: functions.GetVal(l.Protocol),
			Port:     functions.GetVal(l.Port),
		}
		if action.RedirectConfig != nil {
			defaultAction := []DefaultActionConfig{
				{
					RedirectConfig: []RedirectConfig{
						{
							Protocol: functions.GetVal(action.RedirectConfig.Protocol),
							Port:     functions.GetVal(action.RedirectConfig.Port),
						},
					},
				},
			}
			cf.DefaultAction = defaultAction
		}
		resourceConfigs = append(resourceConfigs, AWSResourceConfig{
			Resource: cf,
			Metadata: l.AWSCloudFormationMetadata,
		})
	}

	return resourceConfigs
}
