package accurics

vpcFlowLogEnabled[api.id]{
    api := input.google_compute_subnetwork[_]
    count(api.config.log_config) == 0
}