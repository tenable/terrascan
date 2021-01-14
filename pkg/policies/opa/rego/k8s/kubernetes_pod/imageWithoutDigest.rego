package accurics

#rule for pod, same will satisfy terraform pod, covers containers, initContainers, and terraform init_containers
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.containers[_]
    satisfied := [re_match("@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+", container.image)]
    not all(satisfied)
}

#rule for init containers
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.initContainers[_]
    satisfied := [re_match("@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+", container.image)]
    not all(satisfied)
}

#rule for terraform init_containers
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.init_containers[_]
    satisfied := [re_match("@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+", container.image)]
    not all(satisfied)
}

#rule for deployment, daemonset, job, replica_set, replication_controller, stateful_set covers containers, initContainers, terraform init_containers
{{.prefix}}{{.name}}{{.suffix}}[kind.id] {
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

    kind := item[_]
    container := kind.config.spec.template.spec.containers[_]
    satisfied := [re_match("@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+", container.image)]
    not all(satisfied)
}

{{.prefix}}{{.name}}{{.suffix}}[kind.id] {
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

    kind := item[_]
    container := kind.config.spec.template.spec.initContainers[_]
    satisfied := [re_match("@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+", container.image)]
    not all(satisfied)
}

{{.prefix}}{{.name}}{{.suffix}}[kind.id] {
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

    kind := item[_]
    container := kind.config.spec.template.spec.init_containers[_]
    satisfied := [re_match("@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+", container.image)]
    not all(satisfied)
}

#rule for cron_job, covers containers, initContainers, terraform init_containers
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.containers[_]
    satisfied := [re_match("@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+", container.image)]
    not all(satisfied)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.initContainers[_]
    satisfied := [re_match("@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+", container.image)]
    not all(satisfied)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.init_containers[_]
    satisfied := [re_match("@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+", container.image)]
    not all(satisfied)
}