package accurics

apiGatewayNoAccessLogs[stage.id] {
  stage := input.aws_apigatewayv2_stage[_]
  not stage.config.access_log_settings
}
