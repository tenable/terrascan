package tenable

clusterLabelsEnabled[api.id]{
    api := input.google_container_cluster[_]
    api.config.resource_labels == null
}
