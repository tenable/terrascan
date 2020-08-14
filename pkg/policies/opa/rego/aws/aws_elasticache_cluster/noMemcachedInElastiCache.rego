package accurics

{{.prefix}}noMemcachedInElastiCache[retVal] {
    elasticache = input.aws_elasticache_cluster[_]
    elasticache.config.engine != "redis"
    traverse = "engine"
    retVal := { "Id": elasticache.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "engine", "AttributeDataType": "string", "Expected": "redis", "Actual": elasticache.config.engine }
}
