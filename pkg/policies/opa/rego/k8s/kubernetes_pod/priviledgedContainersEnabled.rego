package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_pod[_]
    pod.config.spec.privileged == true
}

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_pod_security_policy[_]
    pod.config.spec.privileged == true
}