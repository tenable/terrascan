package accurics

{{.prefix}}neptuneClusterLoggingDisabled[neptune_cluster.id] {
    neptune_cluster := input.aws_neptune_cluster[_]
    object.get(neptune_cluster.config, "enable_cloudwatch_logs_exports", "undefined") == "undefined"
}

{{.prefix}}neptuneClusterLoggingDisabled[neptune_cluster.id] {
    neptune_cluster := input.aws_neptune_cluster[_]
    neptune_cluster.config.enable_cloudwatch_logs_exports == []
}

{{.prefix}}neptuneClusterLoggingDisabled[neptune_cluster.id] {
    neptune_cluster := input.aws_neptune_cluster[_]
    neptune_cluster.config.enable_cloudwatch_logs_exports == null
}