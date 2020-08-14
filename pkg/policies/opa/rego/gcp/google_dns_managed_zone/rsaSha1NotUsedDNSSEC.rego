package accurics

rsaSha1NotUsedDNSSEC[api.id]{
    api := input.google_dns_managed_zone[_]
    data := api.config.dnssec_config[_]
    var := data.default_key_specs[_]
    var.algorithm == "rsasha1"
}