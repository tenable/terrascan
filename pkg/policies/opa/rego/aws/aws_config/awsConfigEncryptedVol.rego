package accurics

awsConfigEncryptedVol[api.id] {
    api := input.aws_config_config_rule[_]
    source := api.config.source[_]
    source.source_identifier != "ENCRYPTED_VOLUMES"
}