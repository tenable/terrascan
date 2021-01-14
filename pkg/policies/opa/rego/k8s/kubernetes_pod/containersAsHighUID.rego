package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_pod[_]
    pod.config.spec.securityContext.runAsUser < 1000
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_pod[_]
    container := pod.config.spec.containers[_]
    container.securityContext.runAsUser < 1000
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_pod[_]
    initcontainer := pod.config.spec.initContainers[_]
    initcontainer.securityContext.runAsUser < 1000
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_cron_job[_]
    pod.config.spec.jobTemplate.spec.template.spec.securityContext.runAsUser < 1000
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_cron_job[_]
    container := pod.config.spec.jobTemplate.spec.template.spec.containers[_]
    container.securityContext.runAsUser < 1000
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_cron_job[_]
    initcontainer := pod.config.spec.jobTemplate.spec.template.spec.initContainers[_]
    initcontainer.securityContext.runAsUser < 1000
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_pod_security_policy[_]
    ranges := pod.config.spec.runAsUser.ranges[_]
    ranges.min < 1000
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
    pod.config.spec.template.spec.securityContext.runAsUser < 1000
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
    checkContainer(pod.config.spec.template.spec)
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
    checkInitContainer(pod.config.spec.template.spec)
}

checkInitContainer(spec) {
    containers := spec.initContainers[_]
    containers.securityContext.runAsUser < 1000
}

checkContainer(spec) {
    containers := spec.containers[_]
    containers.securityContext.runAsUser < 1000
}