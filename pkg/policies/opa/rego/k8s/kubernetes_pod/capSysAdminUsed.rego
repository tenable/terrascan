package accurics

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_pod[_]

    some i
    container := pod.config.spec.containers[i]
    container.securityContext.capabilities.add == "-SYS_ADMIN"

    traverse := sprintf("spec.containers[%d].securityContext.capabilities.add", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_pod[_]

    some i
    initcontainer := pod.config.spec.initContainers[i]
    initcontainer.securityContext.capabilities.add == "-SYS_ADMIN"

    traverse := sprintf("spec.initContainers[%d].securityContext.capabilities.add", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_cron_job[_]

    some i
    container := pod.config.spec.jobTemplate.spec.template.spec.containers[i]
    container.securityContext.capabilities.add == "-SYS_ADMIN"

    traverse := sprintf("spec.jobTemplate.spec.template.spec.containers[%d].securityContext.capabilities.add", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_cron_job[_]

    some i
    initcontainer := pod.config.spec.jobTemplate.spec.template.spec.initContainers[i]
    initcontainer.securityContext.capabilities.add == "-SYS_ADMIN"

    traverse := sprintf("spec.jobTemplate.spec.template.spec.initContainers[%d].securityContext.capabilities.add", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
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

    some i
    container := pod.config.spec.template.spec.containers[i]
    container.securityContext.capabilities.add == "-SYS_ADMIN"

    traverse := sprintf("spec.template.spec.containers[%d].securityContext.capabilities.add", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
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

    some i
    container := pod.config.spec.template.spec.initContainers[i]
    container.securityContext.capabilities.add == "-SYS_ADMIN"

    traverse := sprintf("spec.template.spec.initContainers[%d].securityContext.capabilities.add", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}