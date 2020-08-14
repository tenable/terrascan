package accurics

{{.prefix}}keyVaultAuditLoggingEnabled[retVal] {
  kv := input.azurerm_key_vault[_]
  kv.type == "azurerm_key_vault"
  not loggingExist(kv)
  rc = "cmVzb3VyY2UgImF6dXJlcm1fbW9uaXRvcl9kaWFnbm9zdGljX3NldHRpbmciICIjI3Jlc291cmNlX25hbWUjIyIgewogIG5hbWUgICAgICAgICAgICAgICA9ICJhenVyZXJtX2tleV92YXVsdC4jI3Jlc291cmNlX25hbWUjIy5sb2ciCiAgdGFyZ2V0X3Jlc291cmNlX2lkID0gYXp1cmVybV9rZXlfdmF1bHQuIyNyZXNvdXJjZV9uYW1lIyMuaWQKICBzdG9yYWdlX2FjY291bnRfaWQgPSAjI3N0b3JhZ2VfYWNjb3VudF9pZCMjCiAgbG9nIHsKICAgIGNhdGVnb3J5ID0gIkF1ZGl0RXZlbnQiCiAgICBlbmFibGVkICA9IHRydWUKICB9Cn0="
  decode_rc = base64.decode(rc)
  replaced_vpc_id := replace(decode_rc, "##resource_name##", kv.name)
  traverse = ""
  retVal := { "Id": kv.id, "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "resource", "Expected": base64.encode(replaced_vpc_id), "Actual": null }
}

loggingExist(key_vault) = exists {
	log_set := { key_vault_id | key_vault_id := input.azurerm_monitor_diagnostic_setting[i].config.target_resource_id }
	log_set[key_vault.id]
    exists = true
} else = exists {
	log_set := { resource_name | resource_name := input.azurerm_monitor_diagnostic_setting[i].name }
    log_name := sprintf("azurerm_key_vault.%s.log", [key_vault.name])
	log_set[log_name]
    exists = true
}