package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    pod.config.metadata.labels.app == "helm"
    pod.config.metadata.labels.name == "tiller"
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_cron_job[_]
    pod.config.spec.jobTemplate.spec.template.metadata.labels.app == "helm"
    pod.config.spec.jobTemplate.spec.template.metadata.labels.name == "tiller"
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
    checkPod(pod)
}

checkPod(pod) {
    pod.config.spec.template.metadata.labels.app == "helm"
    pod.config.spec.template.metadata.labels.name == "tiller"
}