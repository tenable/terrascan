package accurics

{{.prefix}}checkAuditEnabled[retVal] {
  sql_db_resource := input.azurerm_sql_database[_]
  some i
  threat_detection_policy = sql_db_resource.config.threat_detection_policy[i]
  threat_detection_policy.state == "Disabled"
  
  traverse := sprintf("threat_detection_policy[%d].state", [i])
  retVal := { "Id": sql_db_resource.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "threat_detection_policy.state", "AttributeDataType": "string", "Expected": "Enabled", "Actual": threat_detection_policy.state }
}

{{.prefix}}checkAuditEnabled[retVal] {
  sql_db_resource := input.azurerm_sql_database[_]
  object.get(sql_db_resource.config, "threat_detection_policy", "undefined") == "undefined"
  
  rc := "ewogICJ0aHJlYXRfZGV0ZWN0aW9uX3BvbGljeSI6IHsKICAgICJzdGF0ZSI6ICJFbmFibGVkIiwKICAgICJzdG9yYWdlX2FjY291bnRfYWNjZXNzX2tleSI6ICIke2F6dXJlcm1fc3RvcmFnZV9hY2NvdW50LiMjcmVzb3VyY2VfbmFtZSMjLnByaW1hcnlfYWNjZXNzX2tleX0iLAogICAgInN0b3JhZ2VfZW5kcG9pbnQiOiAiJHthenVyZXJtX3N0b3JhZ2VfYWNjb3VudC4jI3Jlc291cmNlX25hbWUjIy5wcmltYXJ5X2Jsb2JfZW5kcG9pbnR9IiwKICAgICJ1c2Vfc2VydmVyX2RlZmF1bHQiOiAiRW5hYmxlZCIKICB9Cn0="
  decode_rc = base64.decode(rc)
  storage_account := input.azurerm_storage_account[0]
  replaced_resource_name := replace(decode_rc, "##resource_name##", storage_account.name)

  retVal := { "Id": sql_db_resource.id, "ReplaceType": "add", "CodeType": "block", "Traverse": "", "Attribute": "", "AttributeDataType": "block", "Expected": base64.encode(replaced_resource_name), "Actual": null }
}

# create storage_account TODO
# {{.prefix}}checkAuditEnabled[retVal] {
#   sql_db_resource := input.azurerm_sql_database[_]
#   object.get(sql_db_resource.config, "threat_detection_policy", "undefined") == "undefined"
  
#   rc := "ewogICJ0aHJlYXRfZGV0ZWN0aW9uX3BvbGljeSI6IHsKICAgICJzdGF0ZSI6ICJFbmFibGVkIiwKICAgICJzdG9yYWdlX2FjY291bnRfYWNjZXNzX2tleSI6ICIke2F6dXJlcm1fc3RvcmFnZV9hY2NvdW50LiMjcmVzb3VyY2VfbmFtZSMjLnByaW1hcnlfYWNjZXNzX2tleX0iLAogICAgInN0b3JhZ2VfZW5kcG9pbnQiOiAiJHthenVyZXJtX3N0b3JhZ2VfYWNjb3VudC4jI3Jlc291cmNlX25hbWUjIy5wcmltYXJ5X2Jsb2JfZW5kcG9pbnR9IiwKICAgICJ1c2Vfc2VydmVyX2RlZmF1bHQiOiAiRW5hYmxlZCIKICB9Cn0="
#   decode_rc = base64.decode(rc)
#   object.get(input, "azurerm_storage_account", "undefined") == "undefined"
#   replaced_resource_name := replace(decode_rc, "##resource_name##", "blah")

#   retVal := { "Id": sql_db_resource.id, "ReplaceType": "add", "CodeType": "block", "Traverse": "", "Attribute": "", "AttributeDataType": "block", "Expected": base64.encode(replaced_resource_name), "Actual": null }
# }