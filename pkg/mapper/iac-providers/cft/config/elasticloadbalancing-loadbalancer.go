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
	"github.com/awslabs/goformation/v4/cloudformation/elasticloadbalancing"
)

// ElasticLoadBalancingLoadBalancerConfig holds config for aws_elb
type ElasticLoadBalancingLoadBalancerConfig struct {
	Listeners           interface{} `json:"listener"`
	AccessLoggingPolicy interface{} `json:"access_logs,omitempty"`
	Config
}

// ELBAccessLoggingPolicyConfig holds config for access_logs attribute of aws_elb
type ELBAccessLoggingPolicyConfig struct {
	Enabled bool `json:"enabled"`
}

// ELBListenerConfig holds config for listener attribute of aws_elb
type ELBListenerConfig struct {
	LBProtocol       string `json:"lb_protocol"`
	InstanceProtocol string `json:"instance_protocol"`
}

// GetElasticLoadBalancingLoadBalancerConfig returns config for aws_elb
func GetElasticLoadBalancingLoadBalancerConfig(e *elasticloadbalancing.LoadBalancer) []AWSResourceConfig {
	cf := ElasticLoadBalancingLoadBalancerConfig{
		Config: Config{
			Name: e.LoadBalancerName,
			Tags: e.Tags,
		},
	}

	if e.AccessLoggingPolicy != nil {
		cf.AccessLoggingPolicy = ELBAccessLoggingPolicyConfig{
			Enabled: e.AccessLoggingPolicy.Enabled,
		}
	}

	if e.Listeners != nil {
		lc := make([]ELBListenerConfig, 0)
		for _, listener := range e.Listeners {
			lc = append(lc, ELBListenerConfig{
				InstanceProtocol: listener.InstanceProtocol,
				LBProtocol:       listener.Protocol,
			})
		}
		cf.Listeners = lc
	}

	return []AWSResourceConfig{{Resource: cf}}
}
