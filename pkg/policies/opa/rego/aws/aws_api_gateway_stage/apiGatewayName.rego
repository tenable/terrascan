package accurics

apiGatewayName[api.id] {
	api := input.aws_api_gateway_stage[_]
    not checkExists(api.config.stage_name)
}

checkExists(val) = true {
    cloud := input.aws_cloudwatch_log_group[_]
    val == cloud.name
}