package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_cron_job[_]
    checkCorrectAttribute(pod.config.spec.jobTemplate.spec.template.spec)
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
    checkCorrectAttribute(pod.config.spec.template.spec)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    checkCorrectAttribute(pod.config.spec)
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod_security_policy[_]
    podSecurityCheck(pod.config.spec)
}

checkCorrectAttribute(spec) {
    container := spec.containers[_]
    containerSecurityCheck(container)
}

checkCorrectAttribute(spec) {
    container := spec.initContainers[_]
    containerSecurityCheck(container)
}

checkCorrectAttribute(spec) {
    secContext := spec.securityContext
    podSecurityCheck(secContext)
}

containerSecurityCheck(container) {
	{{.not_allowed}}
    container.{{.param1}}.{{.param}} == {{.value}}
}

containerSecurityCheck(container) {
    object.get(container, "{{.param1}}", "undefined") == "undefined"
}

containerSecurityCheck(container) {
    not container.{{.param1}}.{{.param}}
}

containerSecurityCheck(container) {
	{{.allowed}}
    not container.{{.param1}}.{{.arg1}}.{{.arg2}}
}

podSecurityCheck(secContext) {
	{{.not_allowed}}
    secContext.{{.param}} == {{.value}}
}

podSecurityCheck(secContext) {
	{{.not_allowed}}
    object.get(secContext, "{{.param}}", "undefined") == "undefined"
}