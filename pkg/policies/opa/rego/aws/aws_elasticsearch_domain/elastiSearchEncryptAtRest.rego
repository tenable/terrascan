package accurics

elastiSearchEncryptAtRest[api.id] {
    api := input.aws_elasticsearch_domain[_]
    not api.config.encrypt_at_rest
}

elastiSearchEncryptAtRest[api.id] {
    api := input.aws_elasticsearch_domain[_]
    encrypt := api.config.encrypt_at_rest[_]
    encrypt.enabled == false
}

