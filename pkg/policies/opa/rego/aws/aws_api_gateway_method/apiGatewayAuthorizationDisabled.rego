package accurics

{{.prefix}}apiGatewayAuthorizationDisabled[api_gw.id] {
    api_gw := input.aws_api_gateway_method[_]
    api_gw.config.authorization == "NONE"
}