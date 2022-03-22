
package accurics

{{.prefix}}s3BucketAccessLoggingDisabled[retVal] {
    cloudTrail := input.aws_cloudtrail[_]
    bucket := input.aws_s3_bucket[_]

    contains(cloudTrail.config.s3_bucket_name, "{")
    cleanId := cleanSSEBucketID(cloudTrail.config.s3_bucket_name)
    matchBucketName(bucket, cleanId)
    object.get(bucket.config, "logging", "undefined") == ["undefined", [], null][_]
    object.get(input, "aws_s3_bucket_logging", "undefined") == "undefined"

    decode_rc := `resource "aws_s3_bucket_logging" "{{.defaultValue}}" {
        bucket = aws_s3_bucket.##name##.id
        targetBucket = "{{.defaultValue1}}"
    }`

    replaced_rc := replace(decode_rc, "##name##", split(cleanId, ".")[1])

    retVal := {
        "Id": bucket.id,
        "ReplaceType": "add",
        "CodeType": "resource",
        "Traverse": "",
        "Attribute": "",
        "AttributeDataType": "base64",
        "Expected": base64.encode(replaced_rc),
        "Actual": null,
    }
}

{{.prefix}}s3BucketAccessLoggingDisabled[retVal] {
    cloudTrail := input.aws_cloudtrail[_]
    bucket := input.aws_s3_bucket[_]

    not contains(cloudTrail.config.s3_bucket_name, "{")
    matchBucketName(bucket, cloudTrail.config.s3_bucket_name)
    object.get(bucket.config, "logging", "undefined") == ["undefined", [], null][_]
    object.get(input, "aws_s3_bucket_logging", "undefined") == "undefined"

    decode_rc := `resource "aws_s3_bucket_logging" "{{.defaultValue}}" {
        bucket = aws_s3_bucket.##name##.id
        targetBucket = "{{.defaultValue1}}"
    }`

    replaced_rc := replace(decode_rc, "##name##", bucket.name)

    retVal := {
        "Id": bucket.id,
        "ReplaceType": "add",
        "CodeType": "resource",
        "Traverse": "",
        "Attribute": "",
        "AttributeDataType": "base64",
        "Expected": base64.encode(replaced_rc),
        "Actual": null,
    }
}

{{.prefix}}s3BucketAccessLoggingDisabled[retVal] {
    bucket := input.aws_s3_bucket[_]
    object.get(bucket.config, "logging", "undefined") == ["undefined", [], null][_]

    cloudTrail := input.aws_cloudtrail[_]
    not contains(cloudTrail.config.s3_bucket_name, "{")
    matchBucketName(bucket, cloudTrail.config.s3_bucket_name)

    bucket_logged := [bl | bucketlogging := input.aws_s3_bucket_logging[_];
                              bl := bucketlogging.config.bucket]
    not matchBucketName(bucket_logged, bucket)

    decode_rc := `resource "aws_s3_bucket_logging" "{{.defaultValue}}" {
        bucket = aws_s3_bucket.##name##.id
        targetBucket = "{{.defaultValue1}}"
    }`

    replaced_rc := replace(decode_rc, "##name##", bucket.name)

    retVal := {
        "Id": bucket.id,
        "ReplaceType": "add",
        "CodeType": "resource",
        "Traverse": "",
        "Attribute": "",
        "AttributeDataType": "base64",
        "Expected": base64.encode(replaced_rc),
        "Actual": null,
    }
}

{{.prefix}}s3BucketAccessLoggingDisabled[retVal] {
    bucket := input.aws_s3_bucket[_]
    object.get(bucket.config, "logging", "undefined") == ["undefined", [], null][_]

    cloudTrail := input.aws_cloudtrail[_]
    contains(cloudTrail.config.s3_bucket_name, "{")

    cleanId := cleanSSEBucketID(cloudTrail.config.s3_bucket_name)
    matchBucketName(bucket, cleanId)

    bucket_logged := [bl | bucketlogging := input.aws_s3_bucket_logging[_];
                              bl := cleanSSEBucketID(bucketlogging.config.bucket)]
    not matchBucketName(bucket_logged, bucket)

    decode_rc := `resource "aws_s3_bucket_logging" "{{.defaultValue}}" {
        bucket = aws_s3_bucket.##name##.id
        targetBucket = "{{.defaultValue1}}"
    }`

    replaced_rc := replace(decode_rc, "##name##", bucket.name)

    retVal := {
        "Id": bucket.id,
        "ReplaceType": "add",
        "CodeType": "resource",
        "Traverse": "",
        "Attribute": "",
        "AttributeDataType": "base64",
        "Expected": base64.encode(replaced_rc),
        "Actual": null,
    }
}

matchBucketName(arg1, arg2){
    arg1.id == arg2
}

matchBucketName(arg1, arg2){
    arg1[_] == arg2
}

matchBucketName(arg1, arg2){
    arg1[_] == arg2.id
}

matchBucketName(arg1, arg2){
    arg1[_] == arg2.config.id
}

matchBucketName(arg1, arg2){
    arg1[_] == arg2.config.bucket
}

matchBucketName(arg1, arg2){
    arg1.config.id == arg2
}

matchBucketName(arg1, arg2){
    arg1.config.bucket == arg2
}


cleanSSEBucketID(sseBucketID) = cleanID {
    v1 := trim_left(sseBucketID, "$")
    v2 := trim_left(v1, "{")
    v3 := trim_right(v2, "}")
    cleanID = cleanEnd(v3)
}

cleanEnd(sseBucketID_v3) = cleanID {
    endswith(sseBucketID_v3, ".id")
    cleanID = trim_right(sseBucketID_v3, ".id")
}

cleanEnd(sseBucketID_v3) = cleanID {
    endswith(sseBucketID_v3, ".bucket")
    cleanID = trim_right(sseBucketID_v3, ".bucket")
}
