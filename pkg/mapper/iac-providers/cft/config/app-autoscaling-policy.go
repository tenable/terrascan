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
	"github.com/awslabs/goformation/v7/cloudformation/applicationautoscaling"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// AppAutoScalingPolicyConfig holds config for AppAutoScalingPolicy
type AppAutoScalingPolicyConfig struct {
	Config
	Name              string `json:"name"`
	PolicyType        string `json:"policy_type"`
	ResourceID        string `json:"resource_id"`
	ScalableDimension string `json:"scalable_dimension"`
	ServiceNamespace  string `json:"service_namespace"`
}

// GetAppAutoScalingPolicyConfig returns config for AppAutoScalingPolicy
// aws_appautoscaling_policy
func GetAppAutoScalingPolicyConfig(a *applicationautoscaling.ScalingPolicy) []AWSResourceConfig {
	cf := AppAutoScalingPolicyConfig{
		Config: Config{
			Name: a.PolicyName,
		},
		Name:              a.PolicyName,
		PolicyType:        a.PolicyType,
		ResourceID:        functions.GetVal(a.ResourceId),
		ScalableDimension: functions.GetVal(a.ScalableDimension),
		ServiceNamespace:  functions.GetVal(a.ServiceNamespace),
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: a.AWSCloudFormationMetadata,
	}}
}
