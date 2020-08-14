package accurics

{{.prefix}}sqlServerADPredictableAccount[retVal] {
  known_user = { "azure_superuser", "azure_pg_admin", "admin", "administrator", "root", "guest", "public" }
  sql_server := input.azurerm_sql_active_directory_administrator[_]
  sql_server.type == "azurerm_sql_active_directory_administrator"
  user := lower(sql_server.config.login)
  known_user[user]
  uuid_user = uuid.rfc4122(sql_server.config.login)
  traverse := "login"
  retVal := { "Id": sql_server.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "string", "Expected": uuid_user, "Actual": sql_server.config.login }
}
