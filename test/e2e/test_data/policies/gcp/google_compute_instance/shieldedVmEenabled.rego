package accurics

shieldedVmEenabled[api.id] {
     api := input.google_compute_instance[_]
     api.config.shielded_instance_config == []
}

shieldedVmEenabled[api.id] {
     api := input.google_compute_instance[_]
     insConfig := api.config.shielded_instance_config[_]
     insConfig.enable_integrity_monitoring != true
}

shieldedVmEenabled[api.id] {
     api := input.google_compute_instance[_]
     insConfig := api.config.shielded_instance_config[_]
     insConfig.enable_vtpm != true
}