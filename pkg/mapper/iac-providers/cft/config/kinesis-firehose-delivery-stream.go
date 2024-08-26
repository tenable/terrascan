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
	"github.com/awslabs/goformation/v7/cloudformation/kinesisfirehose"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// KinesisFirehoseDeliveryStreamConfig holds config for aws_kinesis_firehose_delivery_stream
type KinesisFirehoseDeliveryStreamConfig struct {
	ServerSideEncryption interface{} `json:"server_side_encryption"`
	Config
}

// KinesisFirehoseDeliveryStreamSseConfig holds config for server_side_encryption attribute of aws_kinesis_firehose_delivery_stream
type KinesisFirehoseDeliveryStreamSseConfig struct {
	KeyType string `json:"key_type,omitempty"`
	KeyARN  string `json:"key_arn,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// GetKinesisFirehoseDeliveryStreamConfig returns aws_kinesis_firehose_delivery_stream
// aws_kinesis_firehose_delivery_stream
func GetKinesisFirehoseDeliveryStreamConfig(k *kinesisfirehose.DeliveryStream) []AWSResourceConfig {
	cf := KinesisFirehoseDeliveryStreamConfig{
		Config: Config{
			Name: functions.GetVal(k.DeliveryStreamName),
			Tags: functions.PatchAWSTags(k.Tags),
		},
	}
	sseConfig := KinesisFirehoseDeliveryStreamSseConfig{}
	if k.DeliveryStreamEncryptionConfigurationInput != nil {
		sseConfig.Enabled = true
		sseConfig.KeyType = k.DeliveryStreamEncryptionConfigurationInput.KeyType
		sseConfig.KeyARN = functions.GetVal(k.DeliveryStreamEncryptionConfigurationInput.KeyARN)
	}
	cf.ServerSideEncryption = []KinesisFirehoseDeliveryStreamSseConfig{sseConfig}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: k.AWSCloudFormationMetadata,
	}}
}
