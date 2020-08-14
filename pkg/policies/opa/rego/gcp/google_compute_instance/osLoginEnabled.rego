package accurics

osLoginEnabled[api.id]
{
    api := input.google_compute_instance[_]
    api.config.metadata != null
    meta_str := lower(object.get(api.config.metadata,"enable-oslogin","undefined"))
    meta_str  == "false"
}

osLoginEnabled[api.id]
{
    api := input.google_compute_project_metadata[_]
    api.config.metadata != null
    meta_str := lower(object.get(api.config.metadata,"enable-oslogin","undefined"))
    meta_str  == "false"    
}