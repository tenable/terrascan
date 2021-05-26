package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/cloudformation"
)

// CloudFormationStackConfig holds config for aws_cloudformation_stack
type CloudFormationStackConfig struct {
	Config
	TemplateURL      interface{} `json:"template_url"`
	NotificationARNs interface{} `json:"notification_arns"`
}

// GetCloudFormationStackConfig returns config for aws_cloudformation_stack
func GetCloudFormationStackConfig(s *cloudformation.Stack) []AWSResourceConfig {
	cf := CloudFormationStackConfig{
		Config: Config{
			Tags: s.Tags,
		},
	}
	if len(s.NotificationARNs) > 0 {
		cf.NotificationARNs = s.NotificationARNs
	}
	if len(s.TemplateURL) > 0 {
		cf.TemplateURL = s.TemplateURL
	}
	return []AWSResourceConfig{{Resource: cf}}
}
