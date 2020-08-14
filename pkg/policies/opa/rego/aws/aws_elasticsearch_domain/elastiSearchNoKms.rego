package accurics

elastiSearchNoKms[api.id]
{
    api := input.aws_elasticsearch_domain[_]
    data := api.config.encrypt_at_rest[_]
    not data.kms_key_id
}

elastiSearchNoKms[api.id]
{
    api := input.aws_elasticsearch_domain[_]
    data := api.config.encrypt_at_rest[_]
    not data.kms_key_id == null
}