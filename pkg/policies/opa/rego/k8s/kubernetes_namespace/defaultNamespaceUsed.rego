package accurics

{{.prefix}}{{.name}}{{.suffix}}[kind.id] {
    item_list := [
        object.get(input, "kubernetes_pod", "undefined"),
        object.get(input, "kubernetes_deployment", "undefined"),
        object.get(input, "kubernetes_job", "undefined"),
    ]

    item = item_list[_]
    item != "undefined"
    kind := item[_]

    checkNamespace(kind.config.metadata)
}

checkNamespace(metadata) {
    lower(metadata.namespace) == "default"
}

checkNamespace(metadata) {
    metadata.namespace == ""
}

checkNamespace(metadata) {
    not metadata.namespace
    not metadata.{{.generate_name}}
}

checkNamespace(metadata) {
    not metadata.namespace
    metadata.{{.generate_name}} == false
}