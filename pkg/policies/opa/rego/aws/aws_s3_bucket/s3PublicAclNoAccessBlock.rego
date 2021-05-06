package accurics

{{.prefix}}s3PublicAclNoAccessBlock[s3_bucket.id] {
    s3_bucket := input.aws_s3_bucket[_]
    startswith(s3_bucket.config.acl, "public")
    object.get(input, "aws_s3_bucket_public_access_block", "undefined") != "undefined"
    public_access_block := input.aws_s3_bucket_public_access_block[_]
    checkPublicAccessBlockExist(s3_bucket.name, public_access_block.config)
}

{{.prefix}}s3PublicAclNoAccessBlock[s3_bucket.id] {
    s3_bucket := input.aws_s3_bucket[_]
    startswith(s3_bucket.config.acl, "public")
    object.get(input, "aws_s3_bucket_public_access_block", "undefined") != "undefined"
}

checkPublicAccessBlockExist(s3_bucket_name, public_access_block) = true {
    bucket_id_split_arr := split(public_access_block.bucket, ".")
    s3_bucket_name == bucket_id_split_arr[1]
}