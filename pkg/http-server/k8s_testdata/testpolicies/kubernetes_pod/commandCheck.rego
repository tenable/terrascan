package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod_kubeapi.id] {
    pod_kubeapi := input.kubernetes_pod[_]
    cmds := pod_kubeapi.config.spec.containers[_].command
    {{.negation}} check(cmds)
}

check(cmds) {
    cmd := cmds[_]
    startswith(cmd, "{{.argument}}")
    {{.presence}} contains(cmd, "{{.param}}")
    {{.optional}}
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    cmds := pod.config.spec.containers[_].imagePullPolicy != "Always"
}