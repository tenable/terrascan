package accurics

#this will satisfy version 4 provider and previous providers
{{.prefix}}noS3BucketSseRules[retVal] {
    bucket := input.aws_s3_bucket[_]
    object.get(bucket.config, "versioning", "undefined") == ["undefined", null, ""][_]
    object.get(input, "aws_s3_bucket_server_side_encryption_configuration", "undefined") == "undefined"
    bucketName := bucket.config.bucket
    decode_rc := `resource "aws_s3_bucket_server_side_encryption_configuration" "example" {
    bucket = $aws_s3_bucket.##name##.bucket$

    rule {
        apply_server_side_encryption_by_default {
        sse_algorithm     = "aws:kms"
        }
    }
    }`

    replacedName := replace(decode_rc, "##name##", bucketName)

    resourceWithoutQuotes := replace(replacedName, "$", " ")

    traverse = ""
    retVal := {"Id": bucket.id, "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "base64", "Expected": base64.encode(resourceWithoutQuotes), "Actual": null}
}

#For terrascan when used without tfplan
{{.prefix}}noS3BucketSseRules[retVal] {
    bucket := input.aws_s3_bucket[_]
    contains(input.aws_s3_bucket_server_side_encryption_configuration[_].config.bucket, "{")
    serverSideEncryptionBucketName := [bName |
        sse := input.aws_s3_bucket_server_side_encryption_configuration[_]
        bName := cleanSSEBucketID(sse.config.bucket)
    ]

    not checkBucketEncrypted(serverSideEncryptionBucketName, bucket)
    decode_rc := `resource "aws_s3_bucket_server_side_encryption_configuration" "example" {
    bucket = $aws_s3_bucket.##name##.bucket$

    rule {
        apply_server_side_encryption_by_default {
        sse_algorithm     = "aws:kms"
        }
    }
    }`

    replacedName := replace(decode_rc, "##name##", bucket.config.bucket)

    resourceWithoutQuotes := replace(replacedName, "$", " ")

    traverse = ""
    retVal := {"Id": bucket.id, "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "base64", "Expected": base64.encode(resourceWithoutQuotes), "Actual": null}
}

checkBucketEncrypted(arg1, arg2) {
   	arg2.id == arg1[_]
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

#For tfplan based scanning
{{.prefix}}noS3BucketSseRules[retVal] {
    bucket := input.aws_s3_bucket[_]
    bucket_name := input.aws_s3_bucket_server_side_encryption_configuration[_].config.bucket
    not contains(bucket_name, "{")
    serverSideEncryptionBucketName := [bName |
        sse := input.aws_s3_bucket_server_side_encryption_configuration[_]
        bName := sse.config.bucket
    ]

    not checkBucketNameForTfPlan(serverSideEncryptionBucketName, bucket)
    decode_rc := `resource "aws_s3_bucket_server_side_encryption_configuration" "example" {
    bucket = $aws_s3_bucket.##name##.bucket$

    rule {
        apply_server_side_encryption_by_default {
        sse_algorithm     = "aws:kms"
        }
    }
    }`

    replacedName := replace(decode_rc, "##name##", bucket.config.bucket)

    resourceWithoutQuotes := replace(replacedName, "$", " ")

    traverse = ""
    retVal := {"Id": bucket.id, "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "base64", "Expected": base64.encode(resourceWithoutQuotes), "Actual": null}
}

checkBucketNameForTfPlan(sse, buck) {
    buck.id == sse[_]
}

checkBucketNameForTfPlan(sse, buck) {
    buck.config.id == sse[_]
}

checkBucketNameForTfPlan(sse, buck) {
    buck.config.bucket == sse[_]
}
