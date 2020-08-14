package accurics

{{.prefix}}apiGatewayLogging[retVal] {
	api := input.aws_api_gateway_stage[_]
	count(api.config.access_log_settings) == 0
	rc = "ewogICJhY2Nlc3NfbG9nX3NldHRpbmdzIjogewogICAgImRlc3RpbmF0aW9uX2FybiI6ICI8ZGVzdGluYXRpb25fYXJuPiIsCiAgICAiZm9ybWF0IjogIjxmb3JtYXQ+PiIKICB9Cn0="
	traverse = ""
	retVal := { "Id": api.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "access_log_settings", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}