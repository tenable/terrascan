package accurics

awsMqPubliclyAccessible[api.id] {
    api := input.aws_mq_broker[_]
    api.config.publicly_accessible == true
}