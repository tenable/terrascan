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
