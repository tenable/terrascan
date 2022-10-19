package accurics

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
    backup := input.aws_db_instance[_]
    object.get(backup.config, "backup_retention_period", "undefined") == ["undefined", null, 0, []][_]

    traverse = "backup_retention_period"
    retVal := { "Id": backup.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "backup_retention_period", "AttributeDataType": "int", "Expected": 30, "Actual": null }
}