package accurics

awsMqLoggingEnabled[api.id] {
    api := input.aws_mq_broker[_]
    api.config.logs == []
}

awsMqLoggingEnabled[api.id] {
    api := input.aws_mq_broker[_]
    propFalse(api.config.logs[_], ["audit", "general"][_])
}

propFalse(obj, prop) = true {
	obj[prop] == false
}