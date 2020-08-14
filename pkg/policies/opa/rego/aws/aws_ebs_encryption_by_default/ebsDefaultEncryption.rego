package accurics

{{.prefix}}ebsDefaultEncryption[retVal] {
  ebsEncrypt := input.aws_ebs_encryption_by_default[_]
  ebsEncrypt.config.enabled == false

  traverse = "enabled"
  retVal := { "Id": ebsEncrypt.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "enabled", "AttributeDataType": "bool", "Expected": true, "Actual": ebsEncrypt.config.enabled }
}