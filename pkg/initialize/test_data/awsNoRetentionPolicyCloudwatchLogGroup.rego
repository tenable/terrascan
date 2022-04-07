package accurics

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
    cw_log_group := input.aws_cloudwatch_log_group[_]
    cw_log_group.config.retention_in_days == 0

    traverse := "retention_in_days"

    retVal := {
        "Id": cw_log_group.id,
        "ReplaceType": "edit",
        "CodeType": "attribute",
        "Traverse": traverse,
        "Attribute": traverse,
        "AttributeDataType": "int",
        "Expected": 120,
        "Actual": cw_log_group.config.retention_in_days
    }
}