package accurics

backupConfigEnabled[api.id] {
    api := input.google_sql_database_instance[_]
    setting := api.config.settings[_]
    backup := setting.backup_configuration[_]
    backup.enabled == false
}