package accurics

awsMqLoggingEnabled[api.id]
{
    api := input.aws_mq_broker[_]
    api.config.logs
}

awsMqLoggingEnabled[api.id]
{
    api := input.aws_mq_broker[_]
    var := api.config.logs[_]
    var.audit == false
}

awsMqLoggingEnabled[api.id]
{
    api := input.aws_mq_broker[_]
    var := api.config.logs[_]
    var.general == false
}
