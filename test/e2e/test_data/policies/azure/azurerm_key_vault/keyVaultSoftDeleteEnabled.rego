package accurics

{{.prefix}}keyVaultSoftDeleteEnabled[retVal] {
  kv := input.azurerm_key_vault[_]
  kv.type == "azurerm_key_vault"
  traverse = "soft_delete_enabled"
  kv.config.soft_delete_enabled == false
  retVal := { "Id": kv.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "bool", "Expected": true, "Actual": kv.config.soft_delete_enabled }
}