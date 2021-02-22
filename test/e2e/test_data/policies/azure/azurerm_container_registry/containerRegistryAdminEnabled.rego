package accurics

{{.prefix}}containerRegistryAdminEnabled[retVal] {
	acr := input.azurerm_container_registry[_]
    acr.config.admin_enabled == true
	traverse = "admin_enabled"
	retVal := { "Id": acr.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "bool", "Expected": false, "Actual": acr.config.admin_enabled }
}