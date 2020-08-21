package accurics

rsaSha1NotUsedDNSSEC[api.id] {
    api := input.google_dns_managed_zone[_]
    dns := api.config.dnssec_config[_]
    keySpec := dns.default_key_specs[_]
    keySpec.algorithm == "rsasha1"
}