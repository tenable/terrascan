package accurics

postgreSqlLogsEnabled[api.id] {
    api := input.azurerm_postgresql_configuration[_]
    api.config.name == "log_checkpoints"
    api.config.value == "off"
}