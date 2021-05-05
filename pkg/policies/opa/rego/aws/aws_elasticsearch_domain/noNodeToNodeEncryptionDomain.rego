package accurics

{{.prefix}}noNodeToNodeEncryptionDomain[domain.id] {
    domain := input.aws_elasticsearch_domain[_]
    object.get(domain.config, "node_to_node_encryption", "undefined") == "undefined"
}

{{.prefix}}noNodeToNodeEncryptionDomain[domain.id] {
    domain := input.aws_elasticsearch_domain[_]
    object.get(domain.config.node_to_node_encryption[_], "enabled", "undefined") == ["undefined", false][_]
}
