package accurics

sslConnectionEnabled[api.id]{
    api := input.azurerm_mysql_server[_]
    not api.config.ssl_enforcement_enabled == true
}