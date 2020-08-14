package accurics

cosNodeImageUsed[api.id]{
    api := input.google_container_node_pool[_]
    data := api.config.node_config[_] 
    not data.image_type == "cos"
}

# cosNodeImageUsed[api.id]{
#   api := input.google_container_node_pool[_]
#   data := api.config.node_config[_] 
#   not data.image_type
#}
