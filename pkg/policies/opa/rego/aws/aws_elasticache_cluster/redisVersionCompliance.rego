package accurics

{{.prefix}}redisVersionCompliance[retVal] {
    elasticache = input.aws_elasticache_cluster[_]
    elasticache.config.engine == "redis"
    engine_version = elasticache.config.engine_version
    min_version_string = "4.0.10"
    min_version = eval_version_number(min_version_string)
    actual_version = eval_version_number(engine_version)
    actual_version < min_version
    traverse = "engine_version"
    retVal := { "Id": elasticache.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "engine_version", "AttributeDataType": "string", "Expected": min_version_string, "Actual": engine_version }
}

eval_version_number(engine_version) = numeric_version {
	version = split(engine_version, ".")
	numeric_version = to_number(version[0]) * 100 + to_number(version[1]) * 10 + to_number(version[2])
}