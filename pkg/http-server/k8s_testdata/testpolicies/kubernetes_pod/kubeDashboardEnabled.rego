package accurics

{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
	pod := input.kubernetes_pod[_]
    pod.config.metadata.labels.app == "kubernetes-dashboard"
}