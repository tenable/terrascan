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
	"fmt"

	"github.com/awslabs/goformation/v7/cloudformation/apigateway"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

const (
	// GatewayMethodSetting represents subresource aws_api_gateway_method_settings for MethodSettings attribute
	GatewayMethodSetting = "MethodSetting"
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
	ClientCertificateID string      `json:"client_certificate_id"`
	RestAPIID           string      `json:"rest_api_id"`
	Config
	XrayTracingEnabled bool `json:"xray_tracing_enabled"`
}

// GetAPIGatewayStageConfig returns config for aws_api_gateway_stage and aws_api_gateway_method_settings
// aws_api_gateway_method_settings
func GetAPIGatewayStageConfig(s *apigateway.Stage) []AWSResourceConfig {

	resourceConfigs := make([]AWSResourceConfig, 0)

	cf := APIGatewayStageConfig{
		Config: Config{
			Name: functions.GetVal(s.StageName),
			Tags: functions.PatchAWSTags(s.Tags),
		},
	}
	cf.RestAPIID = s.RestApiId

	if s.AccessLogSetting != nil {
		cf.AccessLogSettings = s.AccessLogSetting
	} else {
		cf.AccessLogSettings = struct{}{}
	}
	cf.XrayTracingEnabled = functions.GetVal(s.TracingEnabled)
	cf.ClientCertificateID = functions.GetVal(s.ClientCertificateId)

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
					Name: functions.GetVal(s.StageName),
					Tags: functions.PatchAWSTags(s.Tags),
				},
				MethodSettings: []Settings{{
					MetricsEnabled: functions.GetVal(settings.MetricsEnabled),
				}},
			}
			resourceConfigs = append(resourceConfigs, AWSResourceConfig{
				Type: GatewayMethodSetting,
				// Unique name for each method setting used for ID
				Name:     fmt.Sprintf("%s%v", functions.GetVal(s.StageName), i),
				Resource: msc,
				Metadata: s.AWSCloudFormationMetadata,
			})
		}
	}

	return resourceConfigs
}
