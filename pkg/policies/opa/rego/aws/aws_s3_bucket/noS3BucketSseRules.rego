package accurics

{{.prefix}}noS3BucketSseRules[retVal] {
	bucket := input.aws_s3_bucket[_]
    bucket.config.server_side_encryption_configuration == []
    rc = "ewogICJzZXJ2ZXJfc2lkZV9lbmNyeXB0aW9uX2NvbmZpZ3VyYXRpb24iOiB7CiAgICAicnVsZSI6IHsKICAgICAgImFwcGx5X3NlcnZlcl9zaWRlX2VuY3J5cHRpb25fYnlfZGVmYXVsdCI6IHsKICAgICAgICAic3NlX2FsZ29yaXRobSI6ICJBRVMyNTYiCiAgICAgIH0KICAgIH0KICB9Cn0="
    traverse = ""
    retVal := { "Id": bucket.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "server_side_encryption_configuration", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}
