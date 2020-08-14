package accurics

stackDriverMonitoringEnabled[container_cluster.id] {
  container_cluster := input.google_container_cluster[_]
  container_cluster.config.monitoring_service != "monitoring.googleapis.com/kubernetes"
}
