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
	"github.com/awslabs/goformation/v7/cloudformation/dynamodb"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// PITR holds config for point_in_time_recovery block
type PITR struct {
	Enabled bool `json:"enabled"`
}

// SSE holds config for server_side_encryption block
type SSE struct {
	Enabled bool `json:"enabled"`
}

// DynamoDBTableConfig holds config for aws_dynamodb_table
type DynamoDBTableConfig struct {
	Config
	ServerSideEncryption []SSE  `json:"server_side_encryption"`
	PointInTimeRecovery  []PITR `json:"point_in_time_recovery"`
	BillingMode          string `json:"billing_mode"`
}

// GetDynamoDBTableConfig returns config for aws_dynamodb_table
// aws_dynamodb_table
func GetDynamoDBTableConfig(t *dynamodb.Table) []AWSResourceConfig {
	cf := DynamoDBTableConfig{
		Config: Config{
			Tags: t.Tags,
			Name: functions.GetVal(t.TableName),
		},
	}

	if t.BillingMode != nil {
		cf.BillingMode = *t.BillingMode
	}
	if t.SSESpecification != nil {
		cf.ServerSideEncryption = make([]SSE, 1)

		cf.ServerSideEncryption[0].Enabled = t.SSESpecification.SSEEnabled
	}

	if t.PointInTimeRecoverySpecification != nil {
		cf.PointInTimeRecovery = make([]PITR, 1)

		cf.PointInTimeRecovery[0].Enabled = functions.GetVal(t.PointInTimeRecoverySpecification.PointInTimeRecoveryEnabled)
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: t.AWSCloudFormationMetadata,
	}}
}
