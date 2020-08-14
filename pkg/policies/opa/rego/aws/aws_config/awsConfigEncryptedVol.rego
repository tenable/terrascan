package accurics

awsConfigEncryptedVol[api.id]
{
    api := input.aws_config_config_rule[_]
    data := api.config.source[_]
    not data.source_identifier == "ENCRYPTED_VOLUMES"
}