package accurics

sqlAuditingRetention[api.id]{
    api := input.azurerm_sql_server[_]
    var := api.config.extended_auditing_policy[_]
    var.retention_in_days < 90
}