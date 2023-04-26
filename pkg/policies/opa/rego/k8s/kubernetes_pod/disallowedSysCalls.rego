### this policy depends on the parameters specified by the user/client. Here we are considering that no kernel level syscalls are allowed###
package accurics

#rule for pod
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    forbiddenSysctls = ["kernel.*"]
    sysctl := pod.config.spec.securityContext.sysctls[_].name
    forbidden_sysctl(sysctl, forbiddenSysctls)
}

##rule for deployment, daemonset, job, replica_set, stateful_set, replication_controller
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
    forbiddenSysctls = ["kernel.*"]
    sysctl := kind.config.spec.template.spec.securityContext.sysctls[_].name
    forbidden_sysctl(sysctl, forbiddenSysctls)
}

#rule for cron_job
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    forbiddenSysctls = ["kernel.*"]
    sysctl := cron_job.config.spec.jobTemplate.spec.template.spec.securityContext.sysctls[_].name
    forbidden_sysctl(sysctl, forbiddenSysctls)
}

# if all syscalls are forbidden
forbidden_sysctl(sysctl, arg) {
    arg[_] == "*"
}

forbidden_sysctl(sysctl, arg) {
    arg[_] == sysctl
}

forbidden_sysctl(sysctl, arg) {
    startswith(sysctl, trim(arg[_], "*"))
}