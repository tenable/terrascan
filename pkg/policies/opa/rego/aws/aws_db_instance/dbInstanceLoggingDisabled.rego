package accurics

{{.prefix}}dbInstanceLoggingDisabled[db_instance.id] {
    db_instance := input.aws_db_instance[_]
    object.get(db_instance.config, "enabled_cloudwatch_logs_exports", "undefined") == "undefined"
}

{{.prefix}}dbInstanceLoggingDisabled[db_instance.id] {
    db_instance := input.aws_db_instance[_]
    db_instance.config.enabled_cloudwatch_logs_exports == []
}

{{.prefix}}dbInstanceLoggingDisabled[db_instance.id] {
    db_instance := input.aws_db_instance[_]
    db_instance.config.enabled_cloudwatch_logs_exports == null
}