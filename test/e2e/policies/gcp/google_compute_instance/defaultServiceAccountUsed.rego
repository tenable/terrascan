package accurics

defaultServiceAccountUsed[api.id] {
     api := input.google_compute_instance[_]
     fire_rule := api.config.service_account[_]
     contains(fire_rule.email, "@developer.gserviceaccount.com")
}