package accurics

#rule for pod security policy, will be valid for terraform pod_security_policy
{{.prefix}}{{.name}}{{.suffix}}[psp.id] {
    psp := input.kubernetes_pod_security_policy[_]
    psp.config.metadata.annotations["apparmor.security.beta.kubernetes.io/defaultProfileName"] != "runtime/default"
}

#rule for pod, covers containers
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.containers[_]
    not input_apparmor_allowed(container.name, pod.config.metadata)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.initContainers[_]
    not input_apparmor_allowed(container.name, pod.config.metadata)
}

#terraform init_containers
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.init_containers[_]
    not input_apparmor_allowed(container.name, pod.config.metadata)
}

##rule for deployment, daemonset, job, replica_set, stateful_set, replication_controller covers containers
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
    not input_apparmor_allowed(container.name, kind.config.spec.template.metadata)
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
    not input_apparmor_allowed(container.name, kind.config.spec.template.metadata)
}

#terraform init_containers
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
    not input_apparmor_allowed(container.name, kind.config.spec.template.metadata)
}

#rule for cron_job, covers containers
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.containers[_]
    not input_apparmor_allowed(container.name, cron_job.config.spec.jobTemplate.spec.template.metadata)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.initContainers[_]
    not input_apparmor_allowed(container.name, cron_job.config.spec.jobTemplate.spec.template.metadata)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.init_containers[_]
    not input_apparmor_allowed(container.name, cron_job.config.spec.jobTemplate.spec.template.metadata)
}

#function for all Kinds
input_apparmor_allowed(containerName, metadata) {
    metadata.annotations[key] == "runtime/default"
    key == sprintf("container.apparmor.security.beta.kubernetes.io/%v", [containerName])
}