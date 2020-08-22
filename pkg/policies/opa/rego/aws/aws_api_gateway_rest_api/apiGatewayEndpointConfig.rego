package accurics

apiGatewayEndpointConfig[api.id]{
      api := input.aws_api_gateway_rest_api[_]
      endPoint := api.config.endpoint_configuration[_]
      var := endPoint.types[_]
      not var == "PRIVATE"
}