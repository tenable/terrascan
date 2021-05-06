package accurics

{{.prefix}}eksControlPlaneLoggingDisabled[eks_cluster.id] {
    eks_cluster := input.aws_eks_cluster[_]
    object.get(eks_cluster.config, "enabled_cluster_log_types", "undefined") == "undefined"
}

{{.prefix}}eksControlPlaneLoggingDisabled[eks_cluster.id] {
    eks_cluster := input.aws_eks_cluster[_]
    eks_cluster.config.enabled_cluster_log_types == []
}

{{.prefix}}eksControlPlaneLoggingDisabled[eks_cluster.id] {
    eks_cluster := input.aws_eks_cluster[_]
    eks_cluster.config.enabled_cluster_log_types == null
}