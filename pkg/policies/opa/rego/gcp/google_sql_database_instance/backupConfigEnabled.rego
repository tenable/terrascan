package accurics

backupConfigEnabled[api.id]{
    api := input.google_sql_database_instance[_]
    data := api.config.settings[_]
    var := data.backup_configuration[_]
    var.enabled == false
}