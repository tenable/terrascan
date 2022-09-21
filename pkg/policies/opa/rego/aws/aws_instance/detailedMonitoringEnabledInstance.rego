package tenable

{{.prefix}}detailedMonitoringEnabledInstance[instance.id] {
    instance := input.aws_instance[_]
    object.get(instance.config, "monitoring", "undefined") == "undefined"
}

{{.prefix}}detailedMonitoringEnabledInstance[instance.id] {
    instance := input.aws_instance[_]
    instance.config.monitoring != true
}