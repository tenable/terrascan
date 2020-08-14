package accurics

awsAmiEncrypted[api.id]
{
    api := input.aws_ami[_]
    data := api.config.ebs_block_device[_]
    not data.encrypted == true
}