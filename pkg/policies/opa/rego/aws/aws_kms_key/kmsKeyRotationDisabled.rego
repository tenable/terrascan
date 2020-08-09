package accurics

{{.prefix}}kmsKeyRotationDisabled[retVal] {
    kms_key = input.aws_kms_key[_]
    kms_key.config.is_enabled == true
    kms_key.config.enable_key_rotation == false
    traverse = "enable_key_rotation"
    retVal := { "Id": kms_key.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "enable_key_rotation", "AttributeDataType": "bool", "Expected": true, "Actual": kms_key.config.enable_key_rotation }
}