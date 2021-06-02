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
	"encoding/json"

	"github.com/awslabs/goformation/v4/cloudformation/iam"
)

const (
	// IamGroupPolicy represents the sub-resource aws_iam_group_policy for attribute Policy
	IamGroupPolicy = "Policy"
)

// IamGroupPolicyConfig holds config for aws_iam_group_policy
type IamGroupPolicyConfig struct {
	Config
	PolicyName     string `json:"name"`
	PolicyDocument string `json:"policy"`
}

// GetIamGroupConfig returns config for aws_iam_group_policy
func GetIamGroupConfig(r *iam.Group) []AWSResourceConfig {
	// aws_iam_role_policy as a SubResource
	policyConfigs := make([]AWSResourceConfig, 0)
	if r.Policies != nil {
		for _, policy := range r.Policies {
			pc := IamGroupPolicyConfig{
				Config: Config{
					Name: policy.PolicyName,
				},
				PolicyName: policy.PolicyName,
			}
			policyDocument, err := json.Marshal(policy.PolicyDocument)
			if err == nil {
				pc.PolicyDocument = string(policyDocument)
			}
			policyConfigs = append(policyConfigs, AWSResourceConfig{
				Type:     IamGroupPolicy,
				Name:     policy.PolicyName,
				Resource: pc,
			})
		}
	}

	return policyConfigs
}
