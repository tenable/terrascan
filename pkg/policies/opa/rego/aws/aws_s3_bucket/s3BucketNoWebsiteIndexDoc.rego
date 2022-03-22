package accurics

{{.prefix}}s3BucketNoWebsiteIndexDoc[retVal] {
	bucket := input.aws_s3_bucket[_]
    count(bucket.config.website) > 0
    traverse = "website"
    retVal := { "Id": bucket.id, "ReplaceType": "delete", "CodeType": "block", "Traverse": traverse, "Attribute": "website", "AttributeDataType": "block", "Expected": null, "Actual": null }
}

{{.prefix}}s3BucketNoWebsiteIndexDoc[retVal] {
    bucket := input.aws_s3_bucket[_]
    bucket_website_config := input.aws_s3_bucket_website_configuration[_]
    bucket_website_config.config.bucket == bucket.config.bucket
    retVal := { "Id": bucket_website_config.id, "ReplaceType": "delete", "CodeType": "resource", "Traverse": "", "Attribute": "", "AttributeDataType": "resource", "Expected": null, "Actual": null }
}
