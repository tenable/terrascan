package accurics

{{.prefix}}noTags[retVal] {
  cosmos := input.azurerm_cosmosdb_account[_]
  cosmos.config.tags == null
  
  rc := "ewogICJ0YWdzIjogewogICAgImFkZGVkX2J5IjogImFjY3VyaWNzIgogIH0KfQ=="
  retVal := { "Id": cosmos.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": "", "Attribute": "", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}