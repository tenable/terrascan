package accurics

{{.prefix}}{{.name}}{{.suffix}}[dax_cluster.id]{
    dax_cluster := input.aws_dax_cluster[_]
    object.get(dax_cluster.config, "server_side_encryption", "undefined") == ["undefined", []][_]
}

{{.prefix}}{{.name}}{{.suffix}}[dax_cluster.id]{
    dax_cluster := input.aws_dax_cluster[_]
    sseValue := dax_cluster.config.server_side_encryption[_]
    sseValue.enabled == false
}