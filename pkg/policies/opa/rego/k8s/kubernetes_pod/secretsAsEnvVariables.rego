package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_pod[_]
    containers := pod.config.spec.containers[_]
    env := containers.env[_]
    env.valueFrom != []
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_pod[_]
    initcontainer := pod.config.spec.initContainers[_]
    env := initcontainer.env[_]
    env.valueFrom != []
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_cron_job[_]
    containers := pod.config.spec.jobTemplate.spec.template.spec.containers[_]
    env := containers.env[_]
    env.valueFrom != []
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_cron_job[_]
    initcontainer := pod.config.spec.jobTemplate.spec.template.spec.initContainers[_]
    env := initcontainer.env[_]
    env.valueFrom != []
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
    env := containers.env[_]
    env.valueFrom != []
}

checkContainer(spec) {
    containers := spec.containers[_]
    env := containers.env[_]
    env.valueFrom != []
}