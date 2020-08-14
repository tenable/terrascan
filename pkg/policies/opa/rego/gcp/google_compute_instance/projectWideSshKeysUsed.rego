package accurics

projectWideSshKeysUsed[api.id]
{
    api := input.google_compute_instance[_]
    api.config.metadata != null
    meta_str := lower(object.get(api.config.metadata,"block-project-ssh-keys","undefined"))
    meta_str  == "false"
}

projectWideSshKeysUsed[api.id]
{
    api := input.google_compute_project_metadata[_]
    api.config.metadata != null
    meta_str := lower(object.get(api.config.metadata,"block-project-ssh-keys","undefined"))
    meta_str  == "false"    
}