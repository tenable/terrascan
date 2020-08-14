package accurics

{{.prefix}}s3Versioning[retVal] {
	bucket := input.aws_s3_bucket[_]
    some i
    ver := bucket.config.versioning[i]
    ver.enabled == false
    traverse := sprintf("versioning[%d].enabled", [i])
    retVal := { "Id": bucket.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "versioning.enabled", "AttributeDataType": "bool", "Expected": true, "Actual": ver.enabled }
}

{{.prefix}}s3Versioning[retVal] {
	bucket := input.aws_s3_bucket[_]
    not bucket.config.versioning
    rc := "ewogICJ2ZXJzaW9uaW5nIjogewogICAgImVuYWJsZWQiOiB0cnVlCiAgfQp9"
    retVal := { "Id": bucket.id, "ReplaceType": "add", "CodeType": "block", "Attribute": "", "AttributeDataType": "block", "Expected": rc }
}