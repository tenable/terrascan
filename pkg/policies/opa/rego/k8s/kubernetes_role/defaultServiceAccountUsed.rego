package accurics

{{.prefix}}{{.name}}{{.suffix}}[role.id] {
    role := input.kubernetes_cluster_role[_]
    role.config.roleRef.name == "default"
    role.config.roleRef.kind == "role"
}

{{.prefix}}{{.name}}{{.suffix}}[role.id] {
    role := input.kubernetes_role_binding[_]
    role.config.roleRef.name == "default"
    role.config.roleRef.kind == "role"
}