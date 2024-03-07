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
	"github.com/awslabs/goformation/v7/cloudformation/kinesis"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// KinesisStreamConfig holds config for aws_kinesis_stream
type KinesisStreamConfig struct {
	Config
	Name           string `json:"name"`
	KmsKeyID       string `json:"kms_key_id,omitempty"`
	EncryptionType string `json:"encryption_type,omitempty"`
}

// GetKinesisStreamConfig returns config for aws_kinesis_stream
// aws_kinesis_stream
func GetKinesisStreamConfig(k *kinesis.Stream) []AWSResourceConfig {
	cf := KinesisStreamConfig{
		Config: Config{
			Name: functions.GetVal(k.Name),
			Tags: functions.PatchAWSTags(k.Tags),
		},
		Name: functions.GetVal(k.Name),
	}

	if k.StreamEncryption != nil {
		cf.EncryptionType = k.StreamEncryption.EncryptionType
		cf.KmsKeyID = k.StreamEncryption.KeyId
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: k.AWSCloudFormationMetadata,
	}}
}
