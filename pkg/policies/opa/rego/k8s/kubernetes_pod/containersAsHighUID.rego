package accurics

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_pod[_]
    pod.config.spec.securityContext.runAsUser < 1000

    traverse := "spec.securityContext.runAsUser"
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_pod[_]

    some i
    container := pod.config.spec.containers[i]
    container.securityContext.runAsUser < 1000

    traverse := sprintf("spec.containers[%d].securityContext.runAsUser", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_pod[_]

    some i
    initcontainer := pod.config.spec.initContainers[i]
    initcontainer.securityContext.runAsUser < 1000

    traverse := sprintf("spec.containers[%d].securityContext.runAsUser", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_cron_job[_]
    pod.config.spec.jobTemplate.spec.template.spec.securityContext.runAsUser < 1000

    traverse := "spec.jobTemplate.spec.template.spec.securityContext.runAsUser"
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_cron_job[_]

    some i
    container := pod.config.spec.jobTemplate.spec.template.spec.containers[i]
    container.securityContext.runAsUser < 1000

    traverse := sprintf("spec.jobTemplate.spec.template.spec.containers[%d].securityContext.runAsUser", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_cron_job[_]

    some i
    initcontainer := pod.config.spec.jobTemplate.spec.template.spec.initContainers[i]
    initcontainer.securityContext.runAsUser < 1000

    traverse := sprintf("spec.jobTemplate.spec.template.spec.initContainers[%d].securityContext.runAsUser", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
	pod := input.kubernetes_pod_security_policy[_]

    some i
    ranges := pod.config.spec.runAsUser.ranges[i]
    ranges.min < 1000

    traverse := sprintf("spec.runAsUser.ranges[%d].min", [i])
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
    pod.config.spec.template.spec.securityContext.runAsUser < 1000

    traverse := "spec.template.spec.securityContext.runAsUser"
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
    containers.securityContext.runAsUser < 1000

    traverse := sprintf("spec.template.spec.containers[%d].securityContext.runAsUser", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
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
    some i
    container := pod.config.spec.template.spec.initContainers[i]
    containers.securityContext.runAsUser < 1000

    traverse := sprintf("spec.template.spec.initContainers[%d].securityContext.runAsUser", [i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}