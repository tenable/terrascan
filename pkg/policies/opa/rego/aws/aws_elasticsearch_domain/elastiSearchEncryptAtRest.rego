package accurics

elastiSearchEncryptAtRest[api.id]
{
    api := input.aws_elasticsearch_domain[_]
    not api.config.encrypt_at_rest
}

elastiSearchEncryptAtRest[api.id]
{
    api := input.aws_elasticsearch_domain[_]
    data := api.config.encrypt_at_rest[_]
    not data.enabled
}

elastiSearchEncryptAtRest[api.id]
{
    api := input.aws_elasticsearch_domain[_]
    data := api.config.encrypt_at_rest[_]
    data.enabled == false
}

