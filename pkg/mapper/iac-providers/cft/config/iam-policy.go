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
	"encoding/json"

	"github.com/awslabs/goformation/v7/cloudformation/iam"
)

// IamPolicyConfig holds config for aws_iam_policy
type IamPolicyConfig struct {
	Config
	PolicyDocument string `json:"policy"`
	PolicyName     string `json:"name"`
}

// GetIamPolicyConfig returns config for aws_iam_policy
func GetIamPolicyConfig(r *iam.Policy) []AWSResourceConfig {
	cf := IamPolicyConfig{
		Config: Config{
			Name: r.PolicyName,
		},
		PolicyName: r.PolicyName,
	}
	policyDocument, err := json.Marshal(r.PolicyDocument)
	if err == nil {
		cf.PolicyDocument = string(policyDocument)
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: r.AWSCloudFormationMetadata,
	}}
}
