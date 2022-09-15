package accurics

{{.prefix}}{{.name}}{{.suffix}}[service.id] {
    service := input.kubernetes_service[_]
    service_config := service.config
    service_config.spec.type == ["LoadBalancer", "Ingress"][_]
    object.get(service_config.spec, "selector", "undefined") == "undefined"
}