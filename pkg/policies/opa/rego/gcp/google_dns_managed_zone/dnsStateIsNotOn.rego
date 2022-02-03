package accurics

dnsStateIsNotOn[dnsconfig.id] {
  dnsconfig := input.google_dns_managed_zone[_]
  dnssec_config := dnsconfig.config.dnssec_config[_]
  dnssec_config.state != "on"
}
