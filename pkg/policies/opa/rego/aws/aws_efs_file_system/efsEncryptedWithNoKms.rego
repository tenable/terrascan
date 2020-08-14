package accurics

{{.prefix}}{{.name}}[retVal]{
    efs_file = input.aws_efs_file_system[_]
    efs_file.config.encrypted == true 
    not efs_file.config.kms_key_id
    traverse = "kms_key_id"
    retVal := { "Id": efs_file.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "kms_key_id", "AttributeDataType": "string", "Expected": "<kms_key_id>", "Actual": null }
}
