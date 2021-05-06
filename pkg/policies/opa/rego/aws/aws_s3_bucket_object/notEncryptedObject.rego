package accurics

{{.prefix}}notEncryptedObject[s3Object.id] {
    s3Object = input.aws_s3_bucket_object[_]
    object.get(s3Object.config, "kms_key_id", "undefined") == ["", " ", "undefined"][_]
}

{{.prefix}}notEncryptedObject[s3Object.id] {
    s3Object = input.aws_s3_bucket_object[_]
    object.get(s3Object.config, "server_side_encryption", "undefined") == ["", " ", "undefined"][_]
}