package accurics

#rule for pod, covers containers, initContainers, terraform, init_containers
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.containers[_]
    checkForPodLatest(container)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.containers[_]
    checkForPodNoTag(container)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.initContainers[_]
    checkForPodLatest(container)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.initContainers[_]
    checkForPodNoTag(container)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.init_containers[_]
    checkForPodLatest(container)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.init_containers[_]
    checkForPodNoTag(container)
}

#rule for deployment, daemonset, job, replica_set, stateful_set, replication_controller covers containers, initContainers, terraform init_containers
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
    checkForPodLatest(container)
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
    container := kind.config.spec.template.spec.containers[_]
    checkForPodNoTag(container)
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
    checkForPodLatest(container)
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
    checkForPodLatest(container)
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
    checkForPodNoTag(container)
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
    checkForPodNoTag(container)
}

#rule for cron_job, covers containers, initContainers, terraform init_containers
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.containers[_]
    checkForPodLatest(container)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.containers[_]
    checkForPodNoTag(container)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.initContainers[_]
    checkForPodLatest(container)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.init_containers[_]
    checkForPodLatest(container)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.initContainers[_]
    checkForPodNoTag(container)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.init_containers[_]
    checkForPodNoTag(container)
}

#check function for All KINDs
checkForPodLatest(arg) {
    img_split := split(arg.image, ":")
    tag := img_split[count(img_split) - 1]
    tag == "latest"
}

checkForPodNoTag(argument) {
    img_split := split(argument.image, ":")
    count(img_split) == 1
}