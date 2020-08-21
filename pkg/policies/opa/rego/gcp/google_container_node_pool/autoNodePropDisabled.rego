package accurics

{{.name}}[api.id] {
    api := input.google_container_node_pool[_]
    mgmt := api.config.management[_]
    mgmt.{{.property}} == false
}