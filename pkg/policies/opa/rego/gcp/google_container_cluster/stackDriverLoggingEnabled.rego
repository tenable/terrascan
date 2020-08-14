package accurics

stackDriverLoggingEnabled[container_cluster.id] {
  container_cluster := input.google_container_cluster[_]
  container_cluster.config.logging_service != "logging.googleapis.com/kubernetes"
}
