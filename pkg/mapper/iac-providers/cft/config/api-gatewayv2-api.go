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

// CorsConfigurationBlock holds config for cors_configuration attribute
type CorsConfigurationBlock struct {
	AllowCredentials bool     `json:"allow_credentials,omitempty"`
	MaxAge           int      `json:"max_age,omitempty"`
	ExposeHeaders    []string `json:",omitempty"`
	AllowOrigins     []string `json:"allow_origins,omitempty"`
	AllowMethods     []string `json:"allow_methods,omitempty"`
	AllowHeaders     []string `json:"allow_headers,omitempty"`
}

// APIGatewayV2ApiConfig holds config for aws_apigatewayv2_api resource
type APIGatewayV2ApiConfig struct {
	Config
	Name                      string                   `json:"name"`
	ProtocolType              string                   `json:"protocol_type"`
	RouteKey                  string                   `json:"route_key,omitempty"`
	Description               string                   `json:"description,omitempty"`
	CredentialsArn            string                   `json:"credentials_arn,omitempty"`
	RouteSelectionExpression  string                   `json:"route_selection_expression,omitempty"`
	Target                    string                   `json:"target,omitempty"`
	Version                   string                   `json:"version,omitempty"`
	APIKeySelectionExpression string                   `json:"api_key_selection_expression,omitempty"`
	DisableExecuteAPIEndpoint bool                     `json:"disable_execute_api_endpoint,omitempty"`
	FailOnWarnings            bool                     `json:"fail_on_warnings,omitempty"`
	CorsConfiguration         []CorsConfigurationBlock `json:"cors_configuration,omitempty"`
}

// GetAPIGatewayV2ApiConfig returns config for aws_apigatewayv2_api resource
// aws_apigatewayv2_api
func GetAPIGatewayV2ApiConfig(a *apigatewayv2.Api) []AWSResourceConfig {
	var corsConfigData []CorsConfigurationBlock

	if a.CorsConfiguration != nil {
		corsConfigData = make([]CorsConfigurationBlock, 1)
		corsConfigData[0].AllowCredentials = functions.GetVal(a.CorsConfiguration.AllowCredentials)
		corsConfigData[0].AllowHeaders = a.CorsConfiguration.AllowHeaders
		corsConfigData[0].AllowMethods = a.CorsConfiguration.AllowMethods
		corsConfigData[0].AllowOrigins = a.CorsConfiguration.AllowOrigins
		corsConfigData[0].ExposeHeaders = a.CorsConfiguration.ExposeHeaders
		corsConfigData[0].MaxAge = functions.GetVal(a.CorsConfiguration.MaxAge)
	}

	cf := APIGatewayV2ApiConfig{
		Config: Config{
			Name: functions.GetVal(a.Name),
			Tags: functions.PatchAWSTags(a.Tags),
		},
		Name:                      functions.GetVal(a.Name),
		ProtocolType:              functions.GetVal(a.ProtocolType),
		RouteKey:                  functions.GetVal(a.RouteKey),
		Description:               functions.GetVal(a.Description),
		CredentialsArn:            functions.GetVal(a.CredentialsArn),
		RouteSelectionExpression:  functions.GetVal(a.RouteSelectionExpression),
		Target:                    functions.GetVal(a.Target),
		Version:                   functions.GetVal(a.Version),
		APIKeySelectionExpression: functions.GetVal(a.ApiKeySelectionExpression),
		DisableExecuteAPIEndpoint: functions.GetVal(a.DisableExecuteApiEndpoint),
		FailOnWarnings:            functions.GetVal(a.FailOnWarnings),
		CorsConfiguration:         corsConfigData,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: a.AWSCloudFormationMetadata,
	}}
}
