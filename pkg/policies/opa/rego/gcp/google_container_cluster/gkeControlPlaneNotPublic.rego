package accurics

gkeControlPlaneNotPublic[api.id] {
    api := input.google_container_cluster[_]
    pCluster := api.config.private_cluster_config[_]
    pCluster.enable_private_endpoint != true
    object.get(api.config, "master_authorized_networks_config", "undefined") == ["undefined", []][_]
}

gkeControlPlaneNotPublic[api.id] {
    api := input.google_container_cluster[_]
    object.get(api.config, "private_cluster_config", "undefined") == ["undefined", []][_]
    object.get(api.config, "master_authorized_networks_config", "undefined") == ["undefined", []][_]
}