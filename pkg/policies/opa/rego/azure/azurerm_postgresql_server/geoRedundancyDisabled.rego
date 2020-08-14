package accurics

{{.prefix}}geoRedundancyDisabled[retVal] {
  psql_server := input.azurerm_postgresql_server[_]
  psql_server.config.geo_redundant_backup_enabled != true

  traverse = "geo_redundant_backup_enabled"
  retVal := { "Id": psql_server.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "geo_redundant_backup_enabled", "AttributeDataType": "boolean", "Expected": true, "Actual": psql_server.config.geo_redundant_backup_enabled }
}

{{.prefix}}geoRedundancyDisabled[retVal] {
  psql_server := input.azurerm_postgresql_server[_]
  object.get(psql_server.config, "geo_redundant_backup_enabled", "undefined") == "undefined"

  traverse = "geo_redundant_backup_enabled"
  retVal := { "Id": psql_server.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "geo_redundant_backup_enabled", "AttributeDataType": "boolean", "Expected": true, "Actual": null }
}