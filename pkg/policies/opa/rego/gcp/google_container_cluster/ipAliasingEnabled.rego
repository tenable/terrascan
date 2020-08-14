package accurics

ipAliasingEnabled[container_cluster.id] {
  container_cluster := input.google_container_cluster[_]
  container_cluster.config.ip_allocation_policy == []
}
