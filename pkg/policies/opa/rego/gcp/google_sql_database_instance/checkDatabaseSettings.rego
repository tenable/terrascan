package accurics

checkRequireSSLEnabled[db_instance.id] {
  db_instance := input.google_sql_database_instance[_]
  setting := db_instance.config.settings[_]
  not setting.ip_configuration
} {
  db_instance := input.google_sql_database_instance[_]
  setting := db_instance.config.settings[_]  
  ip_configuration = setting.ip_configuration[_]
  not ip_configuration.require_ssl
} {
  db_instance := input.google_sql_database_instance[_]
  setting := db_instance.config.settings[_]  
  ip_configuration = setting.ip_configuration[_]
  ip_configuration.require_ssl == false
}

checkNoPublicAccess[db_instance.id] {
  db_instance := input.google_sql_database_instance[_]
  setting := db_instance.config.settings[_]
  count(setting.ip_configuration) > 0
  ip_configuration = setting.ip_configuration[_]
  count(ip_configuration.authorized_networks) > 0
  authorized_network = ip_configuration.authorized_networks[_]
  authorized_network.value == "0.0.0.0"
}
