package accurics

dnsStateIsNotOn[dnsconfig.id] {
  dnsconfig := input.google_dns_managed_zone[_]
  state := dnsconfig.config.dnssec_config[_].state != "on"
}