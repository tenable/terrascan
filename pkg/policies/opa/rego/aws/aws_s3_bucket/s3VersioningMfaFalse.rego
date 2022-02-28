package accurics

{{.prefix}}s3VersioningMfaFalse[retVal] {
	bucket := input.aws_s3_bucket[_]
    some i
    mfa := bucket.config.versioning[i]
    mfa.mfa_delete == false
    traverse := sprintf("versioning[%d].mfa_delete", [i])
    retVal := { "Id": bucket.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "versioning.mfa_delete", "AttributeDataType": "bool", "Expected": true, "Actual": mfa.mfa_delete }
}

{{.prefix}}s3VersioningMfaFalse[retVal] {
    bucket := input.aws_s3_bucket[_]
    object.get(bucket.config, "versioning", "undefined") == ["undefined", "", null, []][_]
    bucket_versioning := input.aws_s3_bucket_versioning[_]
    some i
    mfa := bucket_versioning.config.versioning_configuration[i]
    lower(mfa.mfa_delete) != "enabled"
    traverse := sprintf("versioning_configuration[%d].mfa_delete", [i])
    retVal := { "Id": bucket_versioning.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "versioning_configuration.mfa_delete", "AttributeDataType": "bool", "Expected": true, "Actual": mfa.mfa_delete }
}