package accurics

{{.prefix}}s3BucketSseRulesWithKmsNull[retVal] {
    bucket := input.aws_s3_bucket[_]
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
	not check_empty(sse.kms_master_key_id)
}

hasEncryption(sse) {
	check_empty(sse.kms_master_key_id)
    sse.sse_algorithm == "AES256"
}
check_empty(key) {
	key == null
}
check_empty(key) {
	key == ""
}
