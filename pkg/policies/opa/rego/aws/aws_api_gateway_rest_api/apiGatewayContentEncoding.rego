package accurics

apiGatewayContentEncoding[api.id]{
    api := input.aws_api_gateway_rest_api[_]
    api.config.minimum_compression_size < 0
}