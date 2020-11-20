package accurics

gkeBasicAuthDisabled[api.id] {
    api := input.google_container_cluster[_]
    auth := api.config.master_auth
    # If username is not specified, basic auth is disabled
    auth[_].username != null

    # If username and password are both empty, basic auth is disabled
    auths := auth[_]
    not gkeBasicAuthEmptyCreds[ auths ]
}

gkeBasicAuthEmptyCreds[auth] {
    auth := input.google_container_cluster[_].config.master_auth[_]
    [auth.username,auth.password] == ["",""]
}
