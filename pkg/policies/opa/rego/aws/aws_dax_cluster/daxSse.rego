package accurics

{{.prefix}}daxSse[dax_cluster.id] {
    dax_cluster := input.aws_dax_cluster[_]
    object.get(dax_cluster.config, "server_side_encryption", "undefined") == [[], "undefined"][_]
}

{{.prefix}}daxSse[dax_cluster.id] {
    dax_cluster := input.aws_dax_cluster[_]
    sse_encryption := dax_cluster.config.server_side_encryption[_]
    sse_encryption.enabled == false
}