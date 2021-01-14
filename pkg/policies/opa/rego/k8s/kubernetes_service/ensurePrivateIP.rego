package accurics

{{.prefix}}{{.name}}{{.suffix}}[service.id] {
    service := input.kubernetes_service[_]
    service.config.kind == "Service"
    type_check(service.config.spec)
    object.get(service.config.spec, "externalIPs", "undefined") != "undefined"
}

type_check(spec) {
    spec.type == "ClusterIP"
}

type_check(spec) {
    object.get(spec, "type", "undefined") == "undefined"
}