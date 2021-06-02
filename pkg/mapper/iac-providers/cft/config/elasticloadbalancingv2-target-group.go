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
	"github.com/awslabs/goformation/v4/cloudformation/elasticloadbalancingv2"
)

// ElasticLoadBalancingV2TargetGroupConfig holds the config for aws_lb_target_group
type ElasticLoadBalancingV2TargetGroupConfig struct {
	Config
	Protocol string `json:"protocol"`
}

// GetElasticLoadBalancingV2TargetGroupConfig returns config for aws_lb_target_group
func GetElasticLoadBalancingV2TargetGroupConfig(l *elasticloadbalancingv2.TargetGroup) []AWSResourceConfig {
	// create a listener subresource per DefaultAction defined in cft
	// as only one default action per listener is possible in terraform
	cf := ElasticLoadBalancingV2TargetGroupConfig{
		Config: Config{
			Name: l.Name,
			Tags: l.Tags,
		},
		Protocol: l.Protocol,
	}

	return []AWSResourceConfig{{Resource: cf}}
}
