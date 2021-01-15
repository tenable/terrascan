package accurics

{{.prefix}}{{.name}}{{.suffix}}[namespace.id] {
    namespace := input.kubernetes_namespace[_]
    object.get(namespace.config.metadata, "labels", "undefined") == "undefined"
}

{{.prefix}}{{.name}}{{.suffix}}[namespace.id] {
    namespace := input.kubernetes_namespace[_]
    object.get(namespace.config.metadata.labels, "owner", "undefined") == "undefined"
}