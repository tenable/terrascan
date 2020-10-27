package accurics

apiGatewayMiconfiguredCors[api.id] {
  api := input.aws_apigatewayv2_api[_]
  cors := api.config.cors_configuration[_]
  origins := cors.allow_origins[_]
  not origins == ["*"]
}
