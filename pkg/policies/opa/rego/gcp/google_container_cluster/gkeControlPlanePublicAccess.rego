package accurics

gkeControlPlanePublicAccess[api.id] {
    api := input.google_container_cluster[_]
    pCluster := api.config.private_cluster_config[_]
    pCluster.enable_private_endpoint != true
    cidr := api.config.master_authorized_networks_config[_].cidr_blocks[_].cidr_block
    checkScopeIsPublic(cidr)
}

scopeIsPrivate(scope) {
    private_ips = ["10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"]
    net.cidr_contains(private_ips[_], scope)
}

checkScopeIsPublic(val) {
    not scopeIsPrivate(val)
}