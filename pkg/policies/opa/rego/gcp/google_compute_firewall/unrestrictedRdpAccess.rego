package accurics

unrestrictedRdpAccess[api.id]
{
     api := input.google_compute_firewall[_]
     data := api.config
     data.direction == "INGRESS"
     fire_rule := data.allow[_]
     fire_rule.protocol == "tcp"
     fire_rule.ports[_] == "3389"
}
