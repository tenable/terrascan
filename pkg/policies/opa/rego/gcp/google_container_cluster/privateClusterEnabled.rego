package accurics

privateClusterEnabled[api.id] {
    api := input.google_container_cluster[_]
    cluster := api.config.private_cluster_config[_]
    cluster.enable_private_endpoint != true
    cluster.enable_private_nodes != true
}