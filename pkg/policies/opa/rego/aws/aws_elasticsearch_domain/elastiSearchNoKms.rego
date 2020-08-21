package accurics

elastiSearchNoKms[api.id] {
    api := input.aws_elasticsearch_domain[_]
    rest := api.config.encrypt_at_rest[_]
    not rest.kms_key_id
}