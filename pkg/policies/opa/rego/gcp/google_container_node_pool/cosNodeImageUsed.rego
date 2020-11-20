package accurics

cosNodeImageUsed[api.id]{
    api := input.google_container_node_pool[_]
    node := api.config.node_config[_]
    not startswith(node.image_type, "cos")
}

