package accurics

{{.prefix}}appGatewayWAFEnabled[retVal] {
	ag := input.azurerm_application_gateway[_]
    ag.type = "azurerm_application_gateway"
	object.get(ag.config, "waf_configuration", "undefined") == "undefined"
    rc = "ewogICJ3YWZfY29uZmlndXJhdGlvbiI6IHsKICAgICJlbmFibGVkIjogdHJ1ZSwKICAgICJmaXJld2FsbF9tb2RlIjogIkRldGVjdGlvbiIsCiAgICAicnVsZV9zZXRfdHlwZSI6ICJPV0FTUCIKICB9Cn0="
    retVal := { "Id": ag.id, "ReplaceType": "add", "CodeType": "block", "Traverse": "", "Attribute": "waf_configuration", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}

{{.prefix}}appGatewayWAFEnabled[retVal] {
	ag := input.azurerm_application_gateway[_]
    ag.type = "azurerm_application_gateway"
	object.get(ag.config, "waf_configuration", "undefined") != "undefined"
	count(ag.config.waf_configuration) <= 0
    rc = "ewogICJ3YWZfY29uZmlndXJhdGlvbiI6IHsKICAgICJlbmFibGVkIjogdHJ1ZSwKICAgICJmaXJld2FsbF9tb2RlIjogIkRldGVjdGlvbiIsCiAgICAicnVsZV9zZXRfdHlwZSI6ICJPV0FTUCIKICB9Cn0="
    retVal := { "Id": ag.id, "ReplaceType": "add", "CodeType": "block", "Traverse": "", "Attribute": "waf_configuration", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}

{{.prefix}}appGatewayWAFEnabled[retVal] {
	ag := input.azurerm_application_gateway[_]
    ag.type = "azurerm_application_gateway"
	object.get(ag.config, "waf_configuration", "undefined") != "undefined"
	count(ag.config.waf_configuration) > 0
    some i
    waf_config := ag.config.waf_configuration[i]
    object.get(waf_config, "enabled", "undefined") == false
    traverse := sprintf("waf_configuration[%d].enabled", [i])
    retVal := { "Id": ag.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "waf_configuration.enabled", "AttributeDataType": "boolean", "Expected": true, "Actual": waf_config.enabled }
}