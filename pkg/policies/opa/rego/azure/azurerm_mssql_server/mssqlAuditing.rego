package accurics

{{.name}}[api.id]{
    api := input.azurerm_mssql_server[_]
    api.config.extended_auditing_policy == []
}

{{.name}}[api.id]{
    api := input.azurerm_mssql_server[_]
    policy := api.config.extended_auditing_policy[_]
    policy.retention_in_days < 90
    {{.checkRetention}}
}