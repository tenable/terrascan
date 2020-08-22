package accurics

apiGatewayTracing[retVal] {
    api := input.aws_api_gateway_stage[_]
    api.config.xray_tracing_enabled == false

    traverse = "xray_tracing_enabled"
    retVal := { "Id": api.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "xray_tracing_enabled", "AttributeDataType": "bool", "Expected": true, "Actual": api.config.xray_tracing_enabled }
}