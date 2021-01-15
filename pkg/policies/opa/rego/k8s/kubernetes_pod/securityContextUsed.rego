package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.containers[_]
    not container.securityContext
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    initcontainer := pod.config.spec.initContainers[_]
    not initcontainer.securityContext
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    not pod.config.spec.securityContext
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_cron_job[_]
    container := pod.config.spec.jobTemplate.spec.template.spec.containers[_]
    not container.securityContext
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_cron_job[_]
    initcontainer := pod.config.spec.jobTemplate.spec.template.spec.initContainers[_]
    not initcontainer.securityContext
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_cron_job[_]
    not pod.config.spec.jobTemplate.spec.template.spec.securityContext
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    item_list := [
        object.get(input, "kubernetes_daemonset", "undefined"),
        object.get(input, "kubernetes_deployment", "undefined"),
        object.get(input, "kubernetes_job", "undefined"),
        object.get(input, "kubernetes_replica_set", "undefined"),
        object.get(input, "kubernetes_replication_controller", "undefined"),
        object.get(input, "kubernetes_stateful_set", "undefined"),
        object.get(input, "kubernetes_cron_job", "undefined")
    ]

    item = item_list[_]
    item != "undefined"

    pod := item[_]
    checkPod(pod)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    item_list := [
        object.get(input, "kubernetes_daemonset", "undefined"),
        object.get(input, "kubernetes_deployment", "undefined"),
        object.get(input, "kubernetes_job", "undefined"),
        object.get(input, "kubernetes_replica_set", "undefined"),
        object.get(input, "kubernetes_replication_controller", "undefined"),
        object.get(input, "kubernetes_stateful_set", "undefined"),
        object.get(input, "kubernetes_cron_job", "undefined")
    ]

    item = item_list[_]
    item != "undefined"

    pod := item[_]
    checkInitContainer(pod.config.spec.template.spec)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    item_list := [
        object.get(input, "kubernetes_daemonset", "undefined"),
        object.get(input, "kubernetes_deployment", "undefined"),
        object.get(input, "kubernetes_job", "undefined"),
        object.get(input, "kubernetes_replica_set", "undefined"),
        object.get(input, "kubernetes_replication_controller", "undefined"),
        object.get(input, "kubernetes_stateful_set", "undefined"),
        object.get(input, "kubernetes_cron_job", "undefined")
    ]

    item = item_list[_]
    item != "undefined"

    pod := item[_]
    checkContainer(pod.config.spec.template.spec)
}

checkContainer(spec) {
    containers := spec.containers[_]
    not containers.securityContext
}

checkInitContainer(spec) {
    containers := spec.initContainers[_]
    not containers.securityContext
}

checkPod(pod) {
    not pod.config.spec.template.spec.securityContext
}