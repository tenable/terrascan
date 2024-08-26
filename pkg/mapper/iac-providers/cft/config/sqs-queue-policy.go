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

	"github.com/awslabs/goformation/v7/cloudformation/sqs"
)

// SqsQueuePolicyConfig holds config for SqsQueuePolicy
type SqsQueuePolicyConfig struct {
	Config
	QueueURL string `json:"queue_url"`
	Policy   string `json:"policy"`
}

// GetSqsQueuePolicyConfig returns config for SqsQueuePolicy
// aws_sqs_queue_policy no policy
func GetSqsQueuePolicyConfig(p *sqs.QueuePolicy) []AWSResourceConfig {
	policyDoc, _ := json.Marshal(p.PolicyDocument)

	cflist := make([]SqsQueuePolicyConfig, len(p.Queues))
	resourcelist := make([]AWSResourceConfig, len(p.Queues))

	for i := range cflist {
		cflist[i].Config.Name = p.Queues[i]
		cflist[i].QueueURL = p.Queues[i]
		cflist[i].Policy = string(policyDoc)

		resourcelist[i].Resource = cflist[i]
		resourcelist[i].Metadata = p.AWSCloudFormationMetadata
	}

	return resourcelist
}
