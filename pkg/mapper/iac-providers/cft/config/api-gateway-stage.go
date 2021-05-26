package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/apigateway"
)

const (
	// GatewayMethodSettings represents subresource aws_api_gateway_method_settings for MethodSettings attribute
	GatewayMethodSettings = "MethodSettings"
)

// MethodSettingConfig holds the config for aws_api_gateway_method_settings
type MethodSettingConfig struct {
	Config
	MetricsEnabled bool `json:"metrics_enabled"`
}

// APIGatewayStageConfig holds config for aws_api_gateway_stage
type APIGatewayStageConfig struct {
	AccessLogSettings   interface{} `json:"access_log_settings"`
	ClientCertificateID interface{} `json:"client_certificate_id"`
	Config
	XrayTracingEnabled bool `json:"xray_tracing_enabled"`
}

// GetAPIGatewayStageConfig returns config for aws_api_gateway_stage and aws_api_gateway_method_settings
func GetAPIGatewayStageConfig(s *apigateway.Stage) []AWSResourceConfig {

	resourceConfigs := make([]AWSResourceConfig, 0)

	cf := APIGatewayStageConfig{
		Config: Config{
			Name: s.StageName,
			Tags: s.Tags,
		},
	}
	if s.AccessLogSetting != nil {
		cf.AccessLogSettings = s.AccessLogSetting
	} else {
		cf.AccessLogSettings = struct{}{}
	}
	cf.XrayTracingEnabled = s.TracingEnabled
	if len(s.ClientCertificateId) > 0 {
		cf.ClientCertificateID = s.ClientCertificateId
	}

	// add aws_api_gateway_stage
	resourceConfigs = append(resourceConfigs, AWSResourceConfig{
		Resource: cf,
	})

	// add aws_api_gateway_method_settings
	// multiple MethodSettings can be configured for same resource in cft
	if s.MethodSettings != nil {
		for _, settings := range s.MethodSettings {
			msc := make(map[string][]MethodSettingConfig)
			msc["settings"] = []MethodSettingConfig{{
				MetricsEnabled: settings.MetricsEnabled,
			}}
			resourceConfigs = append(resourceConfigs, AWSResourceConfig{
				Type:     GatewayMethodSettings,
				Name:     s.StageName,
				Resource: msc,
			})
		}
	}

	return resourceConfigs
}
