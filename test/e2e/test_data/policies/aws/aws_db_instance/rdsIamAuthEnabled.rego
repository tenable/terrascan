package accurics

{{.prefix}}rdsIamAuthEnabled[retVal] {
  rds := input.aws_db_instance[_]
  not rds.config.iam_database_authentication_enabled
  traverse = "iam_database_authentication_enabled"
  retVal := { "Id": rds.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "iam_database_authentication_enabled", "AttributeDataType": "boolean", "Expected": true, "Actual": null }
}

{{.prefix}}rdsIamAuthEnabled[retVal] {
  rds := input.aws_db_instance[_]
  rds.config.iam_database_authentication_enabled == false
  traverse = "iam_database_authentication_enabled"
  retVal := { "Id": rds.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "iam_database_authentication_enabled", "AttributeDataType": "boolean", "Expected": true, "Actual": rds.config.iam_database_authentication_enabled }
}