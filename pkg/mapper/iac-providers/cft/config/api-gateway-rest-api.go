package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/apigateway"
)

// APIGatewayRestAPIConfig holds config for aws_api_gateway_rest_api
type APIGatewayRestAPIConfig struct {
	Config
	EndpointConfiguration  []map[string][]string `json:"endpoint_configuration"`
	MinimumCompressionSize int                   `json:"minimum_compression_size"`
}

// GetAPIGatewayRestAPIConfig returns config for aws_api_gateway_rest_api
func GetAPIGatewayRestAPIConfig(a *apigateway.RestApi) []AWSResourceConfig {
	cf := APIGatewayRestAPIConfig{
		Config: Config{
			Name: a.Name,
			Tags: a.Tags,
		},
		MinimumCompressionSize: a.MinimumCompressionSize,
	}
	// Endpoint Configuration is a []map[string][]string in terraform for some reason
	// despite having fixed keys and not more than one possible value
	ec := make(map[string][]string)
	if a.EndpointConfiguration != nil {
		ec["types"] = a.EndpointConfiguration.Types
		ec["vpc_endpoint_ids"] = a.EndpointConfiguration.VpcEndpointIds
	}
	cf.EndpointConfiguration = []map[string][]string{ec}
	return []AWSResourceConfig{{Resource: cf}}
}
