package accurics

{{.name}}[rule.id] {
  rule := input.google_compute_firewall[_]
  config := rule.config
  config.direction == "INGRESS"
  config.source_ranges[_] == "0.0.0.0/0"
  fire_rule := config.allow[_]
  fire_rule.ports[_] == "{{.port_number}}"
}