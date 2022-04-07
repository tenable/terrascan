package accurics

{{.prefix}}s3BucketSseRulesWithKmsNull[retVal] {
    bucket := input.aws_s3_bucket[_]
    object.get(bucket.config, "server_side_encryption_configuration", "undefined") == ["undefined", null, "", []][_]
    sse := input.aws_s3_bucket_server_side_encryption_configuration[_]
    not contains(sse.config.bucket, "{")
    matchBucketWithSse(sse, bucket)
    retVal := checkCMKKMSUsed(sse)
}

{{.prefix}}s3BucketSseRulesWithKmsNull[retVal] {
    bucket := input.aws_s3_bucket[_]
    object.get(bucket.config, "server_side_encryption_configuration", "undefined") == ["undefined", null, "", []][_]
    sse := input.aws_s3_bucket_server_side_encryption_configuration[_]
    contains(sse.config.bucket, "{")
    cleanId := cleanSSEBucketID(sse.config.bucket)
    matchBucketWithSse(cleanId, bucket)
    retVal := checkCMKKMSUsed(sse)
}

#for AWS provider version < 4
{{.prefix}}s3BucketSseRulesWithKmsNull[retVal] {
    bucket := input.aws_s3_bucket[_]
    object.get(bucket.config, "server_side_encryption_configuration", "undefined") != "undefined"
    object.get(input, "aws_s3_bucket_server_side_encryption_configuration", "undefined") == "undefined"
    some i, j, k
    sse := bucket.config.server_side_encryption_configuration[i]
    sse_rule := sse.rule[j]
    sse_apply := sse_rule.apply_server_side_encryption_by_default[k]

    not hasEncryption(sse_apply)

    sse_apply.kms_master_key_id == null
    traverse := sprintf("server_side_encryption_configuration[%d].rule[%d].apply_server_side_encryption_by_default[%d].kms_master_key_id", [i, j, k])
    retVal := { "Id": bucket.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "server_side_encryption_configuration.rule.apply_server_side_encryption_by_default.kms_master_key_id", "AttributeDataType": "string", "Expected": "<kms_master_key_id>", "Actual": sse_apply.kms_master_key_id }
}

hasEncryption(sse) {
    sse.sse_algorithm == "AES256"
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

matchBucketWithSse(sse, bucket) {
    sse == bucket.name
}

matchBucketWithSse(sse, bucket) {
    sse.config.bucket == bucket.config.bucket
}

matchBucketWithSse(sse, bucket) {
    sse.config.bucket == bucket.config.id
}

matchBucketWithSse(sse, bucket) {
    sse.config.bucket == bucket.id
}

matchBucketWithSse(sse, bucket) {
    sse == bucket.id
}

checkCMKKMSUsed(serverSideEncryption) = retVal {
    object.get(serverSideEncryption.config, "rule", "undefined") == ["undefined", [], null][_]
    decode_rc := `rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = {{.defaultValue}}
      sse_algorithm     = "aws:kms"
    }
      }
    }`

    traverse := ""
    retVal := getRetVal(serverSideEncryption.id, "add", "block", traverse, traverse, "base64", base64.encode(decode_rc), serverSideEncryption.config.rule)
}

checkCMKKMSUsed(serverSideEncryption) = retVal {
    some i
    object.get(serverSideEncryption.config.rule[i], "apply_server_side_encryption_by_default", "undefined") == ["undefined", [], null][_]
    decode_rc := `apply_server_side_encryption_by_default {
      kms_master_key_id = {{.defaultValue}}
      sse_algorithm     = "aws:kms"
    }
    }`

    traverse := sprintf("rule[%d]", [i])
    retVal := getRetVal(serverSideEncryption.id, "add", "block", traverse, "rule", "base64", base64.encode(decode_rc), null)
}

checkCMKKMSUsed(serverSideEncryption) = retVal {
    some i, j
    rules := serverSideEncryption.config.rule[i]
    algoUsed := rules.apply_server_side_encryption_by_default[j].sse_algorithm
    algoUsed == "AES256"
    traverse := sprintf("rule[%d].apply_server_side_encryption_by_default[%d].sse_algorithm", [i, j])
    retVal := getRetVal(serverSideEncryption.id, "edit", "attribute", traverse, "rule.apply_server_side_encryption_by_default.sse_algorithm", "string", "aws:kms", algoUsed)
}

getRetVal(id, ReplaceType, CodeType, Traverse, Attribute, AttributeDataType, Expected, Actual) = retVal {
    retVal := {
        "Id": id,
        "ReplaceType": ReplaceType,
        "CodeType": CodeType,
        "Traverse": Traverse,
        "Attribute": Attribute,
        "AttributeDataType": AttributeDataType,
        "Expected": Expected,
        "Actual": Actual,
    }
}
