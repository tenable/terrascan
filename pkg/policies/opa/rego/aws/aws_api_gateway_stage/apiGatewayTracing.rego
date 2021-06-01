package accurics

apiGatewayTracing[api.id] {
    api := input.aws_api_gateway_stage[_]
    object.get(api.config, "xray_tracing_enabled", "undefined") == [false, "undefined"][_]
}