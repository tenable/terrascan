package accurics

{{.prefix}}lambdaXRayTracingDisabled[retVal] {
    lambda = input.aws_lambda_function[_]
    lambda.type == "aws_lambda_function"
    not lambda.config.tracing_config
    rc = "ewogICJ0cmFjaW5nX2NvbmZpZyI6IHsKICAgICJtb2RlIjogIkFjdGl2ZSIKICB9Cn0="
    retVal := { "Id": lambda.id, "ReplaceType": "add", "CodeType": "block", "Traverse": "", "Attribute": "tracing_config", "AttributeDataType": "base64", "Expected": rc, "Actual": null}
}

{{.prefix}}lambdaXRayTracingDisabled[retVal] {
    lambda = input.aws_lambda_function[_]
    lambda.type == "aws_lambda_function"
    some i
	tracing = lambda.config.tracing_config[i]
    tracing.mode != "Active"
    traverse = sprintf("tracing_config[%d].mode", [i])
    retVal := { "Id": lambda.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "tracing_config.mode", "AttributeDataType": "string", "Expected": "Active", "Actual": tracing.mode }
}