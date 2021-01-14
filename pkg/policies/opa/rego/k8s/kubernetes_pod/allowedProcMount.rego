package accurics

#rule for pod_security_policy
{{.prefix}}{{.name}}{{.suffix}}[psp.id] {
    psp := input.kubernetes_pod_security_policy[_]
    psp.config.spec.allowProcMountTypes != "Default"
}

#rule for pod_security_policy terraform
{{.prefix}}{{.name}}{{.suffix}}[psp.id] {
    psp := input.kubernetes_pod_security_policy[_]
    psp.config.spec.allow_proc_mount_types != "Default"
}

#rule for pod
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    parameters := {}
    container := pod.config.spec.containers[_]
    container.securityContext.procMount
    allowedProcMount := get_allowed_proc_mount(parameters)
    not input_proc_mount_type_allowed(allowedProcMount, container)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    parameters := {}
    container := pod.config.spec.initContainers[_]
    container.securityContext.procMount
    allowedProcMount := get_allowed_proc_mount(parameters)
    not input_proc_mount_type_allowed(allowedProcMount, container)
}

#rule for deployment, daemonset, job, replica_set, stateful_set, replication_controller
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
    #parameters := {  'allowedHostPath' :[{ 'readOnly': true, 'pathPrefix': '/foo' }] }
    parameters := {}
    container.securityContext.procMount
    allowedProcMount := get_allowed_proc_mount(parameters)
    not input_proc_mount_type_allowed(allowedProcMount, container)
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
    #parameters := {  'allowedHostPath' :[{ 'readOnly': true, 'pathPrefix': '/foo' }] }
    parameters := {}
    container.securityContext.procMount
    allowedProcMount := get_allowed_proc_mount(parameters)
    not input_proc_mount_type_allowed(allowedProcMount, container)
}

#rule for cron_job
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.containers[_]
    #parameters := {  'allowedHostPath' :[{ 'readOnly': true, 'pathPrefix': '/foo' }] }
    parameters := {}
    container.securityContext.procMount
    allowedProcMount := get_allowed_proc_mount(parameters)
    not input_proc_mount_type_allowed(allowedProcMount, container)
}

{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    container := cron_job.config.spec.jobTemplate.spec.template.spec.initContainers[_]
    #parameters := {  'allowedHostPath' :[{ 'readOnly': true, 'pathPrefix': '/foo' }] }
    parameters := {}
    container.securityContext.procMount
    allowedProcMount := get_allowed_proc_mount(parameters)
    not input_proc_mount_type_allowed(allowedProcMount, container)
}

###this will get satisfied as no parameters are provided, thus checking with the baseline configuration which is checking that the procmount is default####
get_allowed_proc_mount(params) = out {
    not params.procMount
    out = "default"
}

get_allowed_proc_mount(params) = out {
    not valid_proc_mount(params.procMount)
    out = "default"
}

get_allowed_proc_mount(params) = out {
    out = lower(params.procMount)
}

valid_proc_mount(str) {
    lower(str) == "default"
}

valid_proc_mount(str) {
    lower(str) == "unmasked"
}

input_proc_mount_type_allowed(allowedProcMount, c) {
    allowedProcMount == "default"
    lower(c.securityContext.procMount) == "default"
}