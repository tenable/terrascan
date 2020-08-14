package accurics

autoNodeRepairEnabled[api.id]
{
    api := input.google_container_node_pool[_]
    data := api.config.management[_]
    data.auto_repair == false
}