package accurics

{{.prefix}}{{.name}}{{.suffix}}[api.id]
{
    api := input.{{.resource_type}}[_]
    metadata := api.config.metadata
    metadata.namespace == "default"
}

{{.prefix}}{{.name}}{{.suffix}}[api.id]
{
    api := input.{{.resource_type}}[_]
    metadata := api.config.metadata
    metadata.namespace == ""
}

{{.prefix}}{{.name}}{{.suffix}}[api.id]
{
    api := input.{{.resource_type}}[_]
    metadata := api.config.metadata
    not metadata.namespace
    not metadata.{{.generate_name}}
}

{{.prefix}}{{.name}}{{.suffix}}[api.id]
{
    api := input.{{.resource_type}}[_]
    metadata := api.config.metadata
    not metadata.namespace
    metadata.{{.generate_name}} == false
}

