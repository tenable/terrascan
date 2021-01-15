package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_cron_job[_]
    attribute := pod.config.spec.jobTemplate.spec.template.spec
    boolCheck(attribute)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    attribute := pod.config.spec
    boolCheck(attribute)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    item_list := [
        object.get(input, "kubernetes_daemonset", "undefined"),
        object.get(input, "kubernetes_deployment", "undefined"),
        object.get(input, "kubernetes_job", "undefined"),
        object.get(input, "kubernetes_replica_set", "undefined"),
        object.get(input, "kubernetes_replication_controller", "undefined"),
        object.get(input, "kubernetes_stateful_set", "undefined")
    ]

    item = item_list[_]
    item != "undefined"

    pod := item[_]
    attribute := pod.config.spec.template.spec

    boolCheck(attribute)
}

boolCheck(attribute) {
    attribute.{{.param}} == {{.value}}
}