package accurics

{{.prefix}}{{.name}}{{.suffix}}[service.id] {
    service := input.kubernetes_service[_]
    service_config := service.config
    service_config.spec.type == "NodePort"
}