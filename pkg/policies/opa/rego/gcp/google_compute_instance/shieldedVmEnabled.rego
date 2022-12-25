package accurics

shieldedVmEnabled[api.id] {
     api := input.google_compute_instance[_]
     api.config.shielded_instance_config == []
}

shieldedVmEnabled[api.id] {
     api := input.google_compute_instance[_]
     insConfig := api.config.shielded_instance_config[_]
     insConfig.enable_integrity_monitoring != true
}

shieldedVmEnabled[api.id] {
     api := input.google_compute_instance[_]
     insConfig := api.config.shielded_instance_config[_]
     insConfig.enable_vtpm != true
}