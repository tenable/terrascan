package accurics

defaultServiceAccountUsed[api.id]
{
     api := input.google_compute_instance[_]
     data := api.config
     fire_rule := data.service_account[_]
     contains(fire_rule.email, "@developer.gserviceaccount.com")
}