package accurics

{{.prefix}}rdsAutoMinorVersionUpgradeEnabled[retVal] {
  db := input.aws_db_instance[_]
  db.config.auto_minor_version_upgrade == false
  traverse = "auto_minor_version_upgrade"
  retVal := { "Id": db.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "auto_minor_version_upgrade", "AttributeDataType": "bool", "Expected": true, "Actual": db.config.auto_minor_version_upgrade }
}

