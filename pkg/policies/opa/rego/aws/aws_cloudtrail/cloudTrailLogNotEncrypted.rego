package accurics

{{.prefix}}cloudTrailLogNotEncrypted[retVal]{
    cloud_trail = input.aws_cloudtrail[_]
    cloud_trail.config.kms_key_id == null

    traverse = "kms_key_id"
    retVal := { "Id": cloud_trail.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "kms_key_id", "AttributeDataType": "string", "Expected": "<kms_key_id>", "Actual": cloud_trail.config.kms_key_id }
}