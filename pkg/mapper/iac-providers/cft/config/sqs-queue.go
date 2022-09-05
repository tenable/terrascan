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
	"github.com/awslabs/goformation/v6/cloudformation/sqs"
)

// SqsQueueConfig holds config for SqsQueue
type SqsQueueConfig struct {
	Config
	Name                         string `json:"name"`
	KmsMasterKeyID               string `json:"kms_master_key_id"`
	KmsDataKeyReusePeriodSeconds int    `json:"kms_data_key_reuse_period_seconds"`
	MessageRetentionSeconds      int    `json:"message_retention_seconds"`
}

// GetSqsQueueConfig returns config for SqsQueue
func GetSqsQueueConfig(q *sqs.Queue) []AWSResourceConfig {
	cf := SqsQueueConfig{
		Config: Config{
			Name: *q.QueueName,
		},
		Name:                         *q.QueueName,
		KmsMasterKeyID:               *q.KmsMasterKeyId,
		KmsDataKeyReusePeriodSeconds: *q.KmsDataKeyReusePeriodSeconds,
		MessageRetentionSeconds:      *q.MessageRetentionPeriod,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: q.AWSCloudFormationMetadata,
	}}
}
