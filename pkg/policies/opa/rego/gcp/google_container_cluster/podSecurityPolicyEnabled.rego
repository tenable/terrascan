package accurics

podSecurityPolicyEnabled[api.id]{
    api := input.google_container_cluster[_]
    data := api.config.pod_security_policy_config[_]
    not data.enable_private_endpoint == true
}