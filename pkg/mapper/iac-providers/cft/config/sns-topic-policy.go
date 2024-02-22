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

	"github.com/awslabs/goformation/v7/cloudformation/sns"
)

// SnsTopicPolicyConfig holds config for SnsTopicPolicy
type SnsTopicPolicyConfig struct {
	Config
	ARN    string `json:"arn"`
	Policy string `json:"policy"`
}

// GetSnsTopicPolicyConfig returns config for SnsTopicPolicy
// aws_sns_topic_policy
func GetSnsTopicPolicyConfig(p *sns.TopicPolicy) []AWSResourceConfig {
	policyDoc, _ := json.Marshal(p.PolicyDocument)

	cflist := make([]SnsTopicPolicyConfig, len(p.Topics))
	resourcelist := make([]AWSResourceConfig, len(p.Topics))

	for i := range p.Topics {
		cflist[i].Config.Name = p.Topics[i]
		cflist[i].ARN = p.Topics[i]
		cflist[i].Policy = string(policyDoc)

		resourcelist[i].Resource = cflist[i]
		resourcelist[i].Metadata = p.AWSCloudFormationMetadata
	}

	return resourcelist
}
