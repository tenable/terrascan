package accurics

{{.prefix}}kmsKeyDisabled[retVal] {
    kms_key = input.aws_kms_key[_]
    kms_key.config.is_enabled == false
    traverse = "is_enabled"
    retVal := { "Id": kms_key.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "is_enabled", "AttributeDataType": "bool", "Expected": true, "Actual": kms_key.config.is_enabled }
}