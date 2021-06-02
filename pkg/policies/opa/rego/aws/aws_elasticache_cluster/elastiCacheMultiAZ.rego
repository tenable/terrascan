package accurics

elastiCacheMultiAZ[api.id] {
    api := input.aws_elasticache_cluster[_]
    api.config.az_mode != "cross-az"
}
