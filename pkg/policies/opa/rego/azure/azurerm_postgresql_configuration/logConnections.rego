package accurics

{{.prefix}}logConnections[retVal] {
  psql_config := input.azurerm_postgresql_configuration[_]
  psql_config.config.name == "log_connections"
  psql_config.config.value != "on"

  traverse = "value"
  retVal := { "Id": psql_config.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "value", "AttributeDataType": "string", "Expected": "on", "Actual": psql_config.config.value }
}