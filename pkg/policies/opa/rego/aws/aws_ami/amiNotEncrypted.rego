package accurics

amiNotEncrypted[api.id] {
    api := input.aws_ami[_]
    ebs := api.config.ebs_block_device[_]
    ebs.encrypted != true
}

amiNotEncrypted[api.id] {
    api := input.aws_ami[_]
    object.get(api.config, "ebs_block_device", "undefined") == "undefined"
}