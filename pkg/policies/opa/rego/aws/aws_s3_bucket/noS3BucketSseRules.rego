package accurics

{{.prefix}}noS3BucketSseRules[bucket.id] {
	bucket := input.aws_s3_bucket[_]
    object.get(bucket.config, "server_side_encryption_configuration", "undefined") == [[], null, "undefined"][_]

	sse := input.aws_s3_bucket_server_side_encryption_configuration[_]
    cleanID := cleanSSEBucketID(sse.config.bucket)
    bucket.id == cleanID
    checkSSE(sse.config)
}

# remove all id related prefix and suffix characters generated by terrascan
cleanSSEBucketID(sseBucketID) = cleanID {
    v1 := trim_left(sseBucketID, "$")
    v2 := trim_left(v1, "{")
    v3 := trim_right(v2, "}")
    cleanID := trim_right(v3, ".id")
}

checkSSE(config) {
	object.get(config, "rule", "undefined") == [[], null, "undefined"][_]
}

checkSSE(config) {
	object.get(config.rule[_], "apply_server_side_encryption_by_default", "undefined") == [[], null, "undefined"][_]
}