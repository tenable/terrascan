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
	"github.com/awslabs/goformation/v4/cloudformation/dynamodb"
)

// DynamoDBTableConfig holds config for aws_dynamodb_table
type DynamoDBTableConfig struct {
	Config
	ServerSideEncryption []map[string]interface{} `json:"server_side_encryption"`
}

// GetDynamoDBTableConfig returns config for aws_dynamodb_table
func GetDynamoDBTableConfig(t *dynamodb.Table) []AWSResourceConfig {
	cf := DynamoDBTableConfig{
		Config: Config{
			Tags: t.Tags,
			Name: t.TableName,
		},
	}
	sse := make(map[string]interface{})
	if t.SSESpecification != nil {
		sse["enabled"] = t.SSESpecification.SSEEnabled
	}
	cf.ServerSideEncryption = []map[string]interface{}{sse}
	return []AWSResourceConfig{{Resource: cf}}
}
