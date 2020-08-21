package accurics

gkeControlPlaneNotPublic[api.id] {
    api := input.google_container_cluster[_]
    pCluster := api.config.private_cluster_config[_]
    pCluster.enable_private_endpoint != true
}