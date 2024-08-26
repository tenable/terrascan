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
	"github.com/awslabs/goformation/v7/cloudformation/apigateway"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// APIGatewayRestAPIConfig holds config for aws_api_gateway_rest_api
type APIGatewayRestAPIConfig struct {
	Config
	EndpointConfiguration  []map[string][]string `json:"endpoint_configuration"`
	MinimumCompressionSize int                   `json:"minimum_compression_size"`
	Policy                 interface{}           `json:"policy"`
}

// GetAPIGatewayRestAPIConfig returns config for aws_api_gateway_rest_api
// aws_api_gateway_rest_api
func GetAPIGatewayRestAPIConfig(a *apigateway.RestApi) []AWSResourceConfig {
	cf := APIGatewayRestAPIConfig{
		Config: Config{
			Name: functions.GetVal(a.Name),
			Tags: functions.PatchAWSTags(a.Tags),
		},
		MinimumCompressionSize: functions.GetVal(a.MinimumCompressionSize),
		Policy:                 a.Policy,
	}
	// Endpoint Configuration is a []map[string][]string in terraform for some reason
	// despite having fixed keys and not more than one possible value
	ec := make(map[string][]string)
	if a.EndpointConfiguration != nil {
		ec["types"] = a.EndpointConfiguration.Types
		ec["vpc_endpoint_ids"] = a.EndpointConfiguration.VpcEndpointIds
	}
	cf.EndpointConfiguration = []map[string][]string{ec}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: a.AWSCloudFormationMetadata,
	}}
}
