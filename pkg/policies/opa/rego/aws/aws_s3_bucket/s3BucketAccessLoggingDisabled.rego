package accurics

{{.prefix}}s3BucketAccessLoggingDisabled[s3_bucket.id] {
    s3_bucket := input.aws_s3_bucket[_]
    object.get(s3_bucket.config, "logging", "undefined") == "undefined"
}

{{.prefix}}s3BucketAccessLoggingDisabled[s3_bucket.id] {
    s3_bucket := input.aws_s3_bucket[_]
    s3_bucket.config.logging == []
}

{{.prefix}}s3BucketAccessLoggingDisabled[s3_bucket.id] {
    s3_bucket := input.aws_s3_bucket[_]
    s3_bucket.config.logging == null
}