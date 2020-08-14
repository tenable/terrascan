package accurics

{{.prefix}}{{.name}}[retVal]{
    efs_file = input.aws_efs_file_system[_]
    efs_file.config.encrypted == false
    traverse = "encrypted"
    retVal := { "Id": efs_file.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": efs_file.config.encrypted }
}