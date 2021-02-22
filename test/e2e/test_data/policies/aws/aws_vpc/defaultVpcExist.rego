package accurics

{{.prefix}}defaultVpcExist[vpc.id] {
    vpc = input.aws_vpc[_]
    vpc.config.is_default == true
}