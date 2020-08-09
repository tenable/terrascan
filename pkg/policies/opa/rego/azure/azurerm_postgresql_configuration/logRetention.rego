package accurics

{{.prefix}}logRetention[retVal] {
  psql_config := input.azurerm_postgresql_configuration[_]
  psql_config.config.name == "log_retention_days"
  not checkValid(psql_config.config.value)

  traverse = "value"
  retVal := { "Id": psql_config.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "value", "AttributeDataType": "string", "Expected": "4", "Actual": psql_config.config.value }
}

checkValid(val) = true {
  val == ["4", "5", "6", "7"][_]
}