package accurics

#rule for pod, pod_security_policy covers containers
{{.prefix}}{{.name}}{{.suffix}}[kind.id] {
    item_list := [
        object.get(input, "kubernetes_pod", "undefined"),
        object.get(input, "kubernetes_pod_security_policy", "undefined")
    ]

    item = item_list[_]
    item != "undefined"

    kind := item[_]
    not input_container_allowed(kind.config.metadata)
}

#rule for deployment, daemonset, job, replica_set, stateful_set, replication_controller covers containers
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
    not input_container_allowed(kind.config.spec.template.metadata)
}

#rule for cron_job
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    not input_container_allowed(cron_job.config.spec.jobTemplate.spec.template.metadata)
}

input_container_allowed(metadata) {
    metadata.annotations["seccomp.security.alpha.kubernetes.io/pod"] == "runtime/default"
}

input_container_allowed(metadata) {
    metadata.annotations["seccomp.security.alpha.kubernetes.io/pod"] == "docker/default"
}

  ####Kubernetes v1.19 or later########

#rule for pod covers containers and checks field seccompProfile at container security context which is found at spec.containers.
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.containers[_]
    not check_seccomp(container)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.initContainers[_]
    not check_seccomp(container)
}

#rule for deployment, daemonset, job, replica_set, stateful_set, replication_controller covers containers
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
    not check_seccomp(container)
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
    not check_seccomp(container)
}

#rule for cron_job
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.containers[_]
    not check_seccomp(container)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.initContainers[_]
    not check_seccomp(container)
}

##rule to check seccompProfile at PodSecurityContext which is found at PodSpec##

#rule for pod
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    not check_seccomp(pod.config.spec)
}

#rule for deployment, daemonset, job, replica_set, stateful_set, replication_controller covers containers
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
    not check_seccomp(kind.config.spec.template.spec)
}

#rule for cron_job
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    not check_seccomp(cron_job.config.spec.jobTemplate.spec.template.spec)
}

#function for all Kinds and scenarios
check_seccomp(container) {
    container.securityContext.seccompProfile.type == "RuntimeDefault"
}

check_seccomp(container) {
    container.securityContext.seccompProfile.type == "DockerDefault"
}