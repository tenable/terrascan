package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/ecs"
)

// EcsServiceConfig holds config for aws_ecs_service
type EcsServiceConfig struct {
	Config
	IamRole string `json:"iam_role"`
}

// GetEcsServiceConfig returns config for aws_ecs_service
func GetEcsServiceConfig(c *ecs.Service) []AWSResourceConfig {
	cf := EcsServiceConfig{
		Config: Config{
			Name: c.ServiceName,
			Tags: c.Tags,
		},
		IamRole: c.Role,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
