package accurics

{{.prefix}}{{.name}}[retVal] {
    bucket := input.aws_s3_bucket[_]
    bucket.config.acl == "{{.access}}"
    traverse = "acl"
    retVal := { "Id": bucket.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "acl", "AttributeDataType": "string", "Expected": "private", "Actual": bucket.config.acl }
}