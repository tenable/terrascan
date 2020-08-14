package accurics

legacyAuthEnabled[container_cluster.id] {
  container_cluster := input.google_container_cluster[_]
  container_cluster.config.enable_legacy_abac == true
}
