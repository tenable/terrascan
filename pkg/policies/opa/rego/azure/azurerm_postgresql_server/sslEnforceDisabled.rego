package accurics

{{.prefix}}sslEnforceDisabled[retVal] {
  psql_server := input.azurerm_postgresql_server[_]
  psql_server.config.ssl_enforcement_enabled == false

  traverse = "ssl_enforcement_enabled"
  retVal := { "Id": psql_server.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ssl_enforcement_enabled", "AttributeDataType": "boolean", "Expected": true, "Actual": psql_server.config.ssl_enforcement_enabled }
}
