package accurics

{{.prefix}}storageAccountTrustedMicrosoftServicesEnabled[retVal] {
  storage_account := input.azurerm_storage_account[_]
  some i
  network_rule := storage_account.config.network_rules[i]
  arrayContains(network_rule.bypass, "AzureServices") == false
  traverse = sprintf("network_rules[%d].bypass", [i])
  expected := array.concat(network_rule.bypass, ["AzureServices"])
  retVal := { "Id": storage_account.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "list", "Expected": expected, "Actual": network_rule.bypass }
}

arrayContains(items, elem) = true {
  items[_] = elem
} else = false {
	true
}