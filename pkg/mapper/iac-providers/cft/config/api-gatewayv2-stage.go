package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/apigatewayv2"
)

// APIGatewayV2StageConfig holds config for aws_api_gatewayv2_stage
type APIGatewayV2StageConfig struct {
	Config
	AccessLogSettings interface{} `json:"access_log_settings,omitempty"`
}

// GetAPIGatewayV2StageConfig returns config for aws_api_gatewayv2_stage
func GetAPIGatewayV2StageConfig(s *apigatewayv2.Stage) []AWSResourceConfig {
	cf := APIGatewayV2StageConfig{
		Config: Config{
			Name: s.StageName,
			Tags: s.Tags,
		},
	}
	if s.AccessLogSettings != nil {
		cf.AccessLogSettings = s.AccessLogSettings
	}
	return []AWSResourceConfig{{Resource: cf}}
}
