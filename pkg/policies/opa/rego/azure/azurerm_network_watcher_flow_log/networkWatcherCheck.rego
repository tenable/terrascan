package accurics

{{.prefix}}logRetensionGraterThan90Days[log_object.id] {
  log_object := input.azurerm_network_watcher_flow_log[_]
  retention_policy := log_object.config.retention_policy[_]
  retention_policy.days < 90
}

{{.prefix}}networkWatcherEnabled[log_object.id] {
  log_object := input.azurerm_network_watcher_flow_log[_]
  log_object.config.enabled == false
}
