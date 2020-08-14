package accurics

shieldedVmEenabled[api.id]
{
     api := input.google_compute_instance[_]
     count(api.config.shielded_instance_config) == 0
     
}

shieldedVmEenabled[api.id]
{
     api := input.google_compute_instance[_]
     data := api.config.shielded_instance_config[_]
     not data.enable_integrity_monitoring == true
}

shieldedVmEenabled[api.id]
{
     api := input.google_compute_instance[_]
     data := api.config.shielded_instance_config[_]
     not data.enable_vtpm == true
}