package accurics

elastiCacheMultiAZ[api.id]
{
    api := input.aws_elasticache_cluster[_]
    not api.az_mode == "cross-az"
}