package accurics

{{.prefix}}s3VersioningMfaFalse[retVal] {
	bucket := input.aws_s3_bucket[_]
    some i
    mfa := bucket.config.versioning[i]
    mfa.mfa_delete == false
    traverse := sprintf("versioning[%d].mfa_delete", [i])
    retVal := { "Id": bucket.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "versioning.mfa_delete", "AttributeDataType": "bool", "Expected": true, "Actual": mfa.mfa_delete }
}