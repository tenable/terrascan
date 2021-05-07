package accurics

{{.prefix}}redshiftAuditLogs[redshift.id] {
    redshift := input.aws_redshift_cluster[_]
    object.get(redshift.config, "logging", "undefined") == "undefined"
}

{{.prefix}}redshiftAuditLogs[redshift.id] {
    redshift := input.aws_redshift_cluster[_]
    redshift.config.logging == []
}

{{.prefix}}redshiftAuditLogs[redshift.id] {
    redshift := input.aws_redshift_cluster[_]
    log := redshift.config.logging[_]
    log.enable != true
}