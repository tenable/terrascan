package accurics

{{.prefix}}{{.name}}{{.suffix}}[service.id] {
    service := input.kubernetes_service[_]
    service.config.metadata.labels.name == "tiller-deploy"
}