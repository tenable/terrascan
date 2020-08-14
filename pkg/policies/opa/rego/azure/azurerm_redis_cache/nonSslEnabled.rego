package accurics

{{.prefix}}nonSslEnabled[retVal] {
    redis := input.azurerm_redis_cache[_]
    redis.config.enable_non_ssl_port == true

    traverse = "enable_non_ssl_port"
    retVal := { "Id": redis.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "enable_non_ssl_port", "AttributeDataType": "boolean", "Expected": false, "Actual": redis.config.enable_non_ssl_port }
}