package accurics

{{.prefix}}{{.name}}[retVal] {
    efs_file = input.aws_efs_file_system[_]
    not efs_file.config.encrypted
    traverse = "encrypted"
    retVal := { "Id": efs_file.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "encrypted", "AttributeDataType": "bool", "Expected": true, "Actual": null }
}