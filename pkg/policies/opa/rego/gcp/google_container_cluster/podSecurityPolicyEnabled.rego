package accurics

podSecurityPolicyEnabled[api.id] {
    api := input.google_container_cluster[_]
    policy := api.config.pod_security_policy_config[_]
    policy.enable_private_endpoint != true
}