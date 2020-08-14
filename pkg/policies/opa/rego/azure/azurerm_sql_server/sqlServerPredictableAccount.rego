package accurics

{{.prefix}}sqlServerPredictableAccount[retVal] {
  known_user = { "azure_superuser", "azure_pg_admin", "admin", "administrator", "root", "guest", "public" }
  sql_server := input.azurerm_sql_server[_]
  sql_server.type == "azurerm_sql_server"
  user := lower(sql_server.config.administrator_login)
  known_user[user]
  uuid_user = uuid.rfc4122(sql_server.config.administrator_login)
  traverse := "administrator_login"
  retVal := { "Id": sql_server.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "string", "Expected": uuid_user, "Actual": sql_server.config.administrator_login }
}

{{.prefix}}sqlServerPredictableAccount[retVal] {
  known_user = { "azure_superuser", "azure_pg_admin", "admin", "administrator", "root", "guest", "public" }
  sql_server := input.azurerm_mysql_server[_]
  sql_server.type == "azurerm_mysql_server"
  user := lower(sql_server.config.administrator_login)
  known_user[user]
  uuid_user = uuid.rfc4122(sql_server.config.administrator_login)
  traverse := "administrator_login"
  retVal := { "Id": sql_server.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "string", "Expected": uuid_user, "Actual": sql_server.config.administrator_login }
}


{{.prefix}}sqlServerPredictableAccount[retVal] {
  known_user = { "azure_superuser", "azure_pg_admin", "admin", "administrator", "root", "guest", "public" }
  sql_server := input.azurerm_postgresql_server[_]
  sql_server.type == "azurerm_postgresql_server"
  user := lower(sql_server.config.administrator_login)
  known_user[user]
  uuid_user = uuid.rfc4122(sql_server.config.administrator_login)
  traverse := "administrator_login"
  retVal := { "Id": sql_server.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "string", "Expected": uuid_user, "Actual": sql_server.config.administrator_login }
}

{{.prefix}}sqlServerPredictableAccount[retVal] {
  known_user = { "azure_superuser", "azure_pg_admin", "admin", "administrator", "root", "guest", "public" }
  sql_server := input.azurerm_mssql_server[_]
  sql_server.type == "azurerm_mssql_server"
  user := lower(sql_server.config.administrator_login)
  known_user[user]
  uuid_user = uuid.rfc4122(sql_server.config.administrator_login)
  traverse := "administrator_login"
  retVal := { "Id": sql_server.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "string", "Expected": uuid_user, "Actual": sql_server.config.administrator_login }
}

