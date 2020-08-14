package accurics

apiGatewaySslCertificate[api.id] {
	api := input.aws_api_gateway_stage[_]
    api.config.client_certificate_id == null
}