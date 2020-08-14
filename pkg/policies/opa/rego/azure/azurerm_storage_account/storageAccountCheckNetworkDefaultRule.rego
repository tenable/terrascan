package accurics

{{.prefix}}storageAccountCheckNetworkDefaultRule[retVal] {
  storage_account := input.azurerm_storage_account[_]
  some i
  network_rule := storage_account.config.network_rules[i]
  network_rule.default_action == "Allow"
  traverse = sprintf("network_rules[%d].default_action", [i])
  retVal := { "Id": storage_account.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "string", "Expected": "Deny", "Actual": network_rule.default_action }
}