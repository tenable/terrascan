package accurics

{{.prefix}}atRestNotEncryptedElasticCache[replication_grp.id] {
    replication_grp := input.aws_elasticache_replication_group[_]
    object.get(replication_grp.config, "at_rest_encryption_enabled", "undefined") == [false, "undefined"][_]
}