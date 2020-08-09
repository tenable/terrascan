package accurics

{{.prefix}}checkDataDisksEncrypted[retVal] {
  managed_disk := input.azurerm_managed_disk[_]
  some i
  encryption_settings = managed_disk.config.encryption_settings[i]
  encryption_settings.enabled == false

  traverse := sprintf("encryption_settings[%d].enabled", [i])
  retVal := { "Id": managed_disk.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "encryption_settings.enabled", "AttributeDataType": "boolean", "Expected": true, "Actual": encryption_settings.enabled }
}
