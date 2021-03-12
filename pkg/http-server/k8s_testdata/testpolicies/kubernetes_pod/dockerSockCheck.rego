package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_cron_job[_]
    vol := pod.config.spec.jobTemplate.spec.template.spec.volumes[_]
    socketPathCheck(vol.hostPath.path)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    vol := pod.config.spec.volumes[_]
    socketPathCheck(vol.hostPath.path)
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
    vol := pod.config.spec.template.spec.volumes[_]
    socketPathCheck(vol.hostPath.path)
}

socketPathCheck(attrib) {
    contains(attrib, "/var/run/docker")
}