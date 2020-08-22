package accurics

gkeBasicAuthDisabled[api.id] {
    api := input.google_container_cluster[_]
    auth := api.config.master_auth[_]
    auth.username != null
    auth.password != null
}