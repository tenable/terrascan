package accurics

unrestrictedRdpAccess[api.id] {
     api := input.google_compute_firewall[_]
     api.config.direction == "INGRESS"
     fire_rule := api.config.allow[_]
     fire_rule.protocol == "tcp"
     fire_rule.ports[_] == "3389"
}
