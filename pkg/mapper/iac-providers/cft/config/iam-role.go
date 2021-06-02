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
	// IamRolePolicy represents subresource aws_iam_role_policy for attribute Policy
	IamRolePolicy = "Policy"
)

// IamRoleConfig holds config for aws_iam_role
type IamRoleConfig struct {
	Config
	RoleName                 string `json:"name"`
	AssumeRolePolicyDocument string `json:"assume_role_policy"`
}

// IamRolePolicyConfig holds config for aws_iam_role_policy
type IamRolePolicyConfig struct {
	Config
	PolicyName     string `json:"name"`
	PolicyDocument string `json:"policy"`
}

// GetIamRoleConfig returns config for aws_iam_role and aws_iam_role_policy
func GetIamRoleConfig(r *iam.Role) []AWSResourceConfig {
	resourceConfigs := make([]AWSResourceConfig, 0)

	// add aws_iam_role
	roleConfig := IamRoleConfig{
		Config: Config{
			Name: r.RoleName,
			Tags: r.Tags,
		},
	}
	policyDocument, err := json.Marshal(r.AssumeRolePolicyDocument)
	if err == nil {
		roleConfig.AssumeRolePolicyDocument = string(policyDocument)
	}
	resourceConfigs = append(resourceConfigs, AWSResourceConfig{Resource: roleConfig})

	// aws_iam_role_policy as a SubResource
	// multiple Policies can be defined for a resource in cft
	if r.Policies != nil {
		for _, policy := range r.Policies {
			pc := IamRolePolicyConfig{
				Config: Config{
					Name: policy.PolicyName,
				},
				PolicyName: policy.PolicyName,
			}
			policyDocument, err := json.Marshal(policy.PolicyDocument)
			if err == nil {
				pc.PolicyDocument = string(policyDocument)
			}
			resourceConfigs = append(resourceConfigs, AWSResourceConfig{
				Type:     IamRolePolicy,
				Name:     policy.PolicyName,
				Resource: pc,
			})
		}
	}

	return resourceConfigs
}
