package accurics

sqlServerAuditingEnabled[api.id]{
    api := input.azurerm_sql_server[_]
    count(api.config.extended_auditing_policy) == 0
}