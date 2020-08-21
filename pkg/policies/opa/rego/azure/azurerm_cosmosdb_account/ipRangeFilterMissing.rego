package accurics

{{.prefix}}ipRangeFilterMissing[retVal] {
  cosmos := input.azurerm_cosmosdb_account[_]
  cosmos.config.ip_range_filter == null
  
  traverse := "ip_range_filter"
  retVal := { "Id": cosmos.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ip_range_filter", "AttributeDataType": "string", "Expected": "<cidr>", "Actual": cosmos.config.ip_range_filter }
}