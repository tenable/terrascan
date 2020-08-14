package accurics

masterAuthEnabled[container_cluster.id] {
  container_cluster := input.google_container_cluster[_]
  master := container_cluster.config.master_auth[_]
  master.username == null
}
