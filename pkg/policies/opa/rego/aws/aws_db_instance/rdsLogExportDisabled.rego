package accurics

{{.prefix}}{{.name}}{{.suffix}}[db_instance.id] {
    db_instance := input.aws_db_instance[_]
    object.get(db_instance.config, "enabled_cloudwatch_logs_exports", "undefined") == ["undefined", [], null][_]
}