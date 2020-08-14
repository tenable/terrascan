package accurics

apiGatewayEndpointConfig[api.id]{
      api := input.aws_api_gateway_rest_api[_]
      data := api.config.endpoint_configuration[_]
      var := data.types[_]
      not var == "PRIVATE"
}