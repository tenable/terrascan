package accurics

awsCloudWatchRetentionPreiod[retVal] {
    api := input.aws_cloudwatch_log_group[_]
    api.config.retention_in_days == 0
    
    traverse = "retention_in_days"
    retVal := { "Id": api.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "retention_in_days", "AttributeDataType": "integer", "Expected": "<retention_in_days>", "Actual": api.config.retention_in_days }
}