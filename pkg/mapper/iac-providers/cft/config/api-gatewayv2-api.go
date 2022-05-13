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

import "github.com/awslabs/goformation/v5/cloudformation/apigatewayv2"

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
func GetAPIGatewayV2ApiConfig(a *apigatewayv2.Api) []AWSResourceConfig {
	var corsConfigData []CorsConfigurationBlock

	if a.CorsConfiguration != nil {
		corsConfigData = make([]CorsConfigurationBlock, 1)
		corsConfigData[0].AllowCredentials = a.CorsConfiguration.AllowCredentials
		corsConfigData[0].AllowHeaders = a.CorsConfiguration.AllowHeaders
		corsConfigData[0].AllowMethods = a.CorsConfiguration.AllowMethods
		corsConfigData[0].AllowOrigins = a.CorsConfiguration.AllowOrigins
		corsConfigData[0].ExposeHeaders = a.CorsConfiguration.ExposeHeaders
		corsConfigData[0].MaxAge = a.CorsConfiguration.MaxAge
	}

	cf := APIGatewayV2ApiConfig{
		Config: Config{
			Name: a.Name,
			Tags: a.Tags,
		},
		Name:                      a.Name,
		ProtocolType:              a.ProtocolType,
		RouteKey:                  a.RouteKey,
		Description:               a.Description,
		CredentialsArn:            a.CredentialsArn,
		RouteSelectionExpression:  a.RouteSelectionExpression,
		Target:                    a.Target,
		Version:                   a.Version,
		APIKeySelectionExpression: a.ApiKeySelectionExpression,
		DisableExecuteAPIEndpoint: a.DisableExecuteApiEndpoint,
		FailOnWarnings:            a.FailOnWarnings,
		CorsConfiguration:         corsConfigData,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: a.AWSCloudFormationMetadata,
	}}
}
