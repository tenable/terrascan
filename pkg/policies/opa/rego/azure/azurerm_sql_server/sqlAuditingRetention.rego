package accurics

sqlAuditingRetention[api.id]{
    api := input.azurerm_sql_server[_]
    policy := api.config.extended_auditing_policy[_]
    policy.retention_in_days < 90
}