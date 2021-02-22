package accurics

{{.prefix}}rdsPubliclyAccessible[retVal] {
  db := input.aws_db_instance[_]
  db.config.publicly_accessible == true
  traverse = "publicly_accessible"
  retVal := { "Id": db.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "publicly_accessible", "AttributeDataType": "bool", "Expected": false, "Actual": db.config.publicly_accessible }
}

