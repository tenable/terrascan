package accurics

gkeBasicAuthDisabled[api.id]{
    api := input.google_container_cluster[_]
    data := api.config.master_auth[_]
    not data.username == null
    not data.password == null
}