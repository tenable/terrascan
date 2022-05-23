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
	"fmt"

	"github.com/awslabs/goformation/v5/cloudformation/elasticloadbalancing"
)

// GetPolicies represents subresource aws_load_balancer_policy for Policies attribute
const (
	GetPolicies = "Policies"
)

// PolicyAttributeBlock holds config for PolicyTypeBlock
type PolicyAttributeBlock struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ElasticLoadBalancingLoadBalancerPoliciesConfig holds config for ElasticLoadBalancingLoadBalancerPolicies
type ElasticLoadBalancingLoadBalancerPoliciesConfig struct {
	Config
	LoadBalancerName string                 `json:"load_balancer_name"`
	PolicyName       string                 `json:"policy_name"`
	PolicyTypeName   string                 `jons:"policy_type_name"`
	PolicyAttribute  []PolicyAttributeBlock `json:"policy_attribute"`
}

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
func GetElasticLoadBalancingLoadBalancerConfig(e *elasticloadbalancing.LoadBalancer, elbname string) []AWSResourceConfig {
	elbpolicies := make([]ElasticLoadBalancingLoadBalancerPoliciesConfig, len(e.Policies))
	awsconfig := make([]AWSResourceConfig, len(e.Policies))

	for i := range e.Policies {
		indexedElbName := fmt.Sprintf("%s%d", elbname, i)

		elbpolicies[i].LoadBalancerName = indexedElbName
		elbpolicies[i].PolicyName = e.Policies[i].PolicyName
		elbpolicies[i].PolicyTypeName = e.Policies[i].PolicyType

		elbpolicies[i].PolicyAttribute = make([]PolicyAttributeBlock, len(e.Policies[i].Attributes))
		for ai := range e.Policies[i].Attributes {
			attribVals, ok := e.Policies[i].Attributes[ai].(map[string]interface{})
			if !ok {
				continue
			}

			elbpolicies[i].PolicyAttribute[ai].Name, ok = attribVals["Name"].(string)
			if !ok {
				continue
			}

			elbpolicies[i].PolicyAttribute[ai].Value, ok = attribVals["Value"].(string)
			if !ok {
				continue
			}

			// variable "ok" is only used for safe type conversion
			_ = ok
		}

		awsconfig[i].Type = GetPolicies
		awsconfig[i].Name = indexedElbName
		awsconfig[i].Resource = elbpolicies[i]
		awsconfig[i].Metadata = e.AWSCloudFormationMetadata
	}

	cf := ElasticLoadBalancingLoadBalancerConfig{
		Config: Config{
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

	var awsconfigElb AWSResourceConfig
	awsconfigElb.Resource = cf
	awsconfigElb.Metadata = e.AWSCloudFormationMetadata
	awsconfig = append(awsconfig, awsconfigElb)

	return awsconfig
}
