package accurics

{{.prefix}}dynamoDbEncrypted[dydb_cluster.id] {
    dydb_cluster := input.aws_dynamodb_table[_]
    object.get(dydb_cluster.config, "server_side_encryption", "undefined") == [[], "undefined"][_]
}

{{.prefix}}dynamoDbEncrypted[dydb_cluster.id] {
    dydb_cluster := input.aws_dynamodb_table[_]
    sse_encryption := dydb_cluster.config.server_side_encryption[_]
    sse_encryption.enabled == false
}