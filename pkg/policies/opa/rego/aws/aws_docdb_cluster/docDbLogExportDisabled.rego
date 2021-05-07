package accurics

{{.prefix}}docDbLogExportDisabled[docdb.id] {
    docdb := input.aws_docdb_cluster[_]
    object.get(docdb.config, "enabled_cloudwatch_logs_exports", "undefined") == "undefined"
}

{{.prefix}}docDbLogExportDisabled[docdb.id] {
    docdb := input.aws_docdb_cluster[_]
    docdb.config.enabled_cloudwatch_logs_exports == []
}

{{.prefix}}docDbLogExportDisabled[docdb.id] {
    docdb := input.aws_docdb_cluster[_]
    docdb.config.enabled_cloudwatch_logs_exports == null
}