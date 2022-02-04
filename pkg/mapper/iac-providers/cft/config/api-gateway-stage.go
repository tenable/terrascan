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
	"fmt"

	"github.com/awslabs/goformation/v5/cloudformation/apigateway"
)

const (
	// GatewayMethodSettings represents subresource aws_api_gateway_method_settings for MethodSettings attribute
	GatewayMethodSettings = "MethodSettings"
)

// MethodSettingConfig holds config for aws_api_gateway_method_settings
type MethodSettingConfig struct {
	Config
	MethodSettings []Settings `json:"settings"`
}

// Settings holds configs for the MethodSetting attribute
type Settings struct {
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
		Metadata: s.AWSCloudFormationMetadata,
	})

	// add aws_api_gateway_method_settings
	// multiple MethodSettings can be configured for same resource in cft
	if s.MethodSettings != nil {
		for i, settings := range s.MethodSettings {
			msc := MethodSettingConfig{
				Config: Config{
					Name: s.StageName,
					Tags: s.Tags,
				},
				MethodSettings: []Settings{{
					MetricsEnabled: settings.MetricsEnabled,
				}},
			}
			resourceConfigs = append(resourceConfigs, AWSResourceConfig{
				Type: GatewayMethodSettings,
				// Unique name for each method setting used fopr ID
				Name:     fmt.Sprintf("%s%v", s.StageName, i),
				Resource: msc,
				Metadata: s.AWSCloudFormationMetadata,
			})
		}
	}

	return resourceConfigs
}
