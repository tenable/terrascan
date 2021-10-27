package accurics

{{.prefix}}{{.name}}{{.suffix}}[deployment.id] {
    deployment := input.kubernetes_deployment[_]
    image := deployment.config.spec.template.spec.containers[_].image

    contains(image, "ingress-nginx/controller")
    contains(image, "@sha")
    version := split(split(image, ":v")[1], "@")
    isVulnerableVersion(version)
    isAllowSnippetAnnotations(deployment.config.metadata.namespace)
}

{{.prefix}}{{.name}}{{.suffix}}[deployment.id] {
    deployment := input.kubernetes_deployment[_]
    image := deployment.config.spec.template.spec.containers[_].image

    contains(image, "ingress-nginx/controller")
    not contains(image, "@sha")
    version := split(image, ":v")
    isVulnerableVersion(version)
    isAllowSnippetAnnotations(deployment.metadata.namespace)
}

{{.prefix}}{{.name}}{{.suffix}}[deployment.id] {
    deployment := input.kubernetes_deployment[_]
    image := deployment.config.spec.template.spec.containers[_].image

    contains(image, "ingress-nginx/controller")
    contains(image, "@sha")
    version := split(split(image, ":v")[1], "@")
    isVulnerableVersion(version)
    isAllowSnippetAnnotations(deployment.config.metadata.namespace)
}

{{.prefix}}{{.name}}{{.suffix}}[deployment.id] {
    deployment := input.kubernetes_deployment[_]
    image := deployment.config.spec.template.spec.containers[_].image

    contains(image, "ingress-nginx/controller")
    not contains(image, "@sha")
    version := split(image, ":v")
    isVulnerableVersion(version)
    isAllowSnippetAnnotations(deployment.metadata.namespace)
}

isVulnerableVersion(ver) {
    ver[minus(count(ver), 1)] <= "0.49"
}

isVulnerableVersion(ver) {
    ver[minus(count(ver), 1)] == "1.0.0"
}

isVulnerableVersion(ver) {
    ver[0] <= "0.49"
}

isVulnerableVersion(ver) {
    ver[0] == "1.0.0"
}

isAllowSnippetAnnotations(namespace) {
    configmap := input.kubernetes_config_map[_]
    configmap.config.metadata.namespace == namespace
    configmap.config.data["allow-snippet-annotations"] == "true"
}
