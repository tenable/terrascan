package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod_security_policy[_]
    pod.config.spec.{{.attribute}} != []
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.containers[_]
    container.{{.attribute}} != []
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    initcontainer := pod.config.spec.initContainers[_]
    initcontainer.{{.attribute}} != []
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_cron_job[_]
    container := pod.config.spec.jobTemplate.spec.template.spec.containers[_]
    container.{{.attribute}} != []
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_cron_job[_]
    initcontainer := pod.config.spec.jobTemplate.spec.template.spec.initContainers[_]
    initcontainer.{{.attribute}} != []
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
    containerCheck(pod.config.spec.template.spec)
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
    initContainerCheck(pod.config.spec.template.spec)
}

initContainerCheck(spec) {
    container := spec.initContainers[_]
    container.{{.attribute}} != []
}

containerCheck(spec) {
    container := spec.containers[_]
    container.{{.attribute}} != []
}