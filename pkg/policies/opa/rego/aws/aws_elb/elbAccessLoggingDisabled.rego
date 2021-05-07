package accurics

{{.prefix}}elbAccessLoggingDisabled[elb.id] {
    elb := input.aws_elb[_]
    object.get(elb.config, "access_logs", "undefined") == "undefined"
}

{{.prefix}}elbAccessLoggingDisabled[elb.id] {
    elb := input.aws_elb[_]
    elb.config.access_logs == []
}

{{.prefix}}elbAccessLoggingDisabled[elb.id] {
    elb := input.aws_elb[_]
    access_logs := elb.config.access_logs[_]
    access_logs.enabled != true
}