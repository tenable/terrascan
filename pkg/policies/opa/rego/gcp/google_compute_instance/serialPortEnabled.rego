package accurics

serialPortEnabled[api.id]
{
    api := input.google_compute_instance[_]
    api.config.metadata != null
    meta_str := lower(object.get(api.config.metadata,"serial-port-enable","undefined"))
    meta_str  == "false"
}

serialPortEnabled[api.id]
{
    api := input.google_compute_project_metadata[_]
    api.config.metadata != null
    meta_str := lower(object.get(api.config.metadata,"serial-port-enable","undefined"))
    meta_str  == "false"    
}
