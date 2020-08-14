package accurics

{{.prefix}}kmsKeyNoDeletionWindow[retVal] {
    kms_key = input.aws_kms_key[_]
    kms_key.config.is_enabled == true
    kms_key.config.enable_key_rotation == true
    invalid_window_in_days(kms_key.config.deletion_window_in_days) == true
    traverse = "deletion_window_in_days"
    retVal := { "Id": kms_key.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "deletion_window_in_days", "AttributeDataType": "int", "Expected": 90, "Actual": kms_key.config.deletion_window_in_days }
}

invalid_window_in_days(days) = true {
    days == null
}

invalid_window_in_days(days) = true {
    days != null
    days > 90
}