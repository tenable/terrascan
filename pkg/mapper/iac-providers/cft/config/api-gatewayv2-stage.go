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
	"github.com/awslabs/goformation/v7/cloudformation/apigatewayv2"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// APIGatewayV2StageConfig holds config for aws_api_gatewayv2_stage
type APIGatewayV2StageConfig struct {
	Config
	AccessLogSettings interface{} `json:"access_log_settings,omitempty"`
}

// GetAPIGatewayV2StageConfig returns config for aws_api_gatewayv2_stage
// aws_api_gatewayv2_stage no policy
func GetAPIGatewayV2StageConfig(s *apigatewayv2.Stage) []AWSResourceConfig {
	cf := APIGatewayV2StageConfig{
		Config: Config{
			Name: s.StageName,
			Tags: functions.PatchAWSTags(s.Tags),
		},
	}
	if s.AccessLogSettings != nil {
		cf.AccessLogSettings = s.AccessLogSettings
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: s.AWSCloudFormationMetadata,
	}}
}
