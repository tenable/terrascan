package accurics

{{.prefix}}storageAccountEnableHttps[retVal] {
  enablehttp := input.azurerm_storage_account[_]
  enablehttp.config.enable_https_traffic_only == false
  traverse := "enable_https_traffic_only"
  retVal := { "Id": enablehttp.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "string", "Expected": true, "Actual": enablehttp.config.enable_https_traffic_only }
}