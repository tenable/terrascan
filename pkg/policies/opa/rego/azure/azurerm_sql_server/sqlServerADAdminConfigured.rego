package accurics

{{.prefix}}sqlServerADAdminConfigured[retVal] {
  sql_server := input.azurerm_sql_server[_]
  sql_server.type == "azurerm_sql_server"
  key := concat("-", [sql_server.config.resource_group_name, sql_server.config.name])
  not adAdminExist(key)
  rc = "ZGF0YSAiYXp1cmVybV9jbGllbnRfY29uZmlnIiAiY3VycmVudCIge30KCnJlc291cmNlICJhenVyZXJtX3NxbF9hY3RpdmVfZGlyZWN0b3J5X2FkbWluaXN0cmF0b3IiICIjI3Jlc291cmNlX25hbWUjIyIgewogIHNlcnZlcl9uYW1lICAgICAgICAgPSBhenVyZXJtX3NxbF9zZXJ2ZXIuIyNyZXNvdXJjZV9uYW1lIyMubmFtZQogIHJlc291cmNlX2dyb3VwX25hbWUgPSBhenVyZXJtX3Jlc291cmNlX2dyb3VwLiMjcmVzb3VyY2VfbmFtZSMjLm5hbWUKICBsb2dpbiAgICAgICAgICAgICAgID0gInNxbGFkbWluIgogIHRlbmFudF9pZCAgICAgICAgICAgPSBkYXRhLmF6dXJlcm1fY2xpZW50X2NvbmZpZy5jdXJyZW50LnRlbmFudF9pZAogIG9iamVjdF9pZCAgICAgICAgICAgPSBkYXRhLmF6dXJlcm1fY2xpZW50X2NvbmZpZy5jdXJyZW50Lm9iamVjdF9pZAp9"
  decode_rc = base64.decode(rc)
  replaced := replace(decode_rc, "##resource_name##", sql_server.name)
  traverse = ""
  retVal := { "Id": sql_server.id, "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "resource", "Expected": base64.encode(replaced), "Actual": null }
}

adAdminExist(rg_servername) = exists {
	ad_admin_set := { ad_id | input.azurerm_sql_active_directory_administrator[i].type == "azurerm_sql_active_directory_administrator"; ad_id := concat("-", [input.azurerm_sql_active_directory_administrator[i].config.resource_group_name, input.azurerm_sql_active_directory_administrator[i].config.server_name]) }
	ad_admin_set[rg_servername]
    exists = true
} else = false {
	true
}