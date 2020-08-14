package accurics

networkPolicyEnabled[container_cluster.id] {
  container_cluster := input.google_container_cluster[_]
  container_cluster.config.network_policy[_].enabled == false
}
