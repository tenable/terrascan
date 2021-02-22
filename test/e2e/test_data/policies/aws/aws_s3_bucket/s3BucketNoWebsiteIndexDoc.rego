package accurics

{{.prefix}}s3BucketNoWebsiteIndexDoc[retVal] {
	bucket := input.aws_s3_bucket[_]
    count(bucket.config.website) > 0
    traverse = "website"
    retVal := { "Id": bucket.id, "ReplaceType": "delete", "CodeType": "block", "Traverse": traverse, "Attribute": "website", "AttributeDataType": "block", "Expected": null, "Actual": null }
}