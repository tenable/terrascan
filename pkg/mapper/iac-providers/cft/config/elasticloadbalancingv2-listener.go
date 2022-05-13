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
	"github.com/awslabs/goformation/v5/cloudformation/elasticloadbalancingv2"
)

// ElasticLoadBalancingV2ListenerConfig holds config for aws_lb_listener
type ElasticLoadBalancingV2ListenerConfig struct {
	Config
	Protocol      string              `json:"protocol"`
	DefaultAction DefaultActionConfig `json:"default_action"`
}

// DefaultActionConfig holds config for default_action attribute of aws_lb_listener
type DefaultActionConfig struct {
	RedirectConfig RedirectConfig `json:"redirect"`
}

// RedirectConfig holds config for redirect attirbute of default_action
type RedirectConfig struct {
	Protocol string `json:"protocol"`
}

// GetElasticLoadBalancingV2ListenerConfig returns config for aws_lb_listener
func GetElasticLoadBalancingV2ListenerConfig(l *elasticloadbalancingv2.Listener) []AWSResourceConfig {
	// create a listener subresource per DefaultAction defined in cft
	// as only one default action per listener is possible in terraform
	resourceConfigs := make([]AWSResourceConfig, 0)

	for _, action := range l.DefaultActions {
		// DefaultActions are required
		cf := ElasticLoadBalancingV2ListenerConfig{
			Config:   Config{},
			Protocol: l.Protocol,
		}
		if action.RedirectConfig != nil {
			defaultAction := DefaultActionConfig{
				RedirectConfig: RedirectConfig{
					Protocol: action.RedirectConfig.Protocol,
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
