package accurics

#when resource is not declared, this is for every s3 bucket
{{.prefix}}s3Versioning[retVal] {
    s3Bucket := input.aws_s3_bucket[_]
    object.get(s3Bucket.config, "versioning", "undefined") == ["undefined", null, ""][_]
    object.get(input, "aws_s3_bucket_versioning", "undefined") == "undefined"
    bucketName := s3Bucket.config.bucket

    decode_rc := `resource "aws_s3_bucket_versioning" "example" {
    bucket = "aws_s3_bucket.##name##.bucket"

    versioning_configuration {
        status = "Enabled"
    }
    }`

    replacedName := replace(decode_rc, "##name##", bucketName)

    resourceWithoutQuotes := replace(replacedName, "\"", " ")

    traverse = ""
    retVal := getRetVal(s3Bucket.id, "add", "resource", traverse, "", "base64", base64.encode(resourceWithoutQuotes), null)
}

#to satisfy previous version of aws provider
{{.prefix}}s3Versioning[retVal] {
    s3Bucket := input.aws_s3_bucket[_]
    some i
    versionData := s3Bucket.config.versioning[i]
    versionData.enabled == false

    traverse := sprintf("versioning[%d].enabled", [i])
    retVal := getRetVal(s3Bucket.id, "edit", "attribute", traverse, "versioning.enabled", "bool", true, versionData.enabled)
}

#To satisfy specific s3 buckets which do not have versioning configuration resource for terrascan without tfplan
{{.prefix}}s3Versioning[retVal] {
    s3Bucket := input.aws_s3_bucket[_]
    object.get(input, "aws_s3_bucket_versioning", "undefined") != "undefined"
    contains(input.aws_s3_bucket_versioning[_].config.bucket, "{")
    bucketVersioning := [bv | bucketVersion := input.aws_s3_bucket_versioning[_];
                              bv := cleanSSEBucketID(bucketVersion.config.bucket)]
    not checkBucketName(s3Bucket, bucketVersioning)

    bucketName := s3Bucket.config.bucket

    decode_rc := `resource "aws_s3_bucket_versioning" "example" {
    bucket = "aws_s3_bucket.##name##.bucket"

    versioning_configuration {
        status = "Enabled"
    }
    }`

    replacedName := replace(decode_rc, "##name##", bucketName)

    resourceWithoutQuotes := replace(replacedName, "\"", " ")

    traverse = ""
    retVal := getRetVal(s3Bucket.id, "add", "resource", traverse, "", "base64", base64.encode(resourceWithoutQuotes), null)
}

checkBucketName(arg1, arg2){
    arg1.id == arg2[_]
}

#To satisfy specific s3 buckets which do not have versioning configuration resource with tfplan
{{.prefix}}s3Versioning[retVal] {
    s3Bucket := input.aws_s3_bucket[_]
    object.get(input, "aws_s3_bucket_versioning", "undefined") != "undefined"
    bucket_version := input.aws_s3_bucket_versioning[_].config.bucket
    not contains(bucket_version, "{")
    bucketVersioning := [bv | bucketVersion := input.aws_s3_bucket_versioning[_];
                              bv := bucketVersion.config.bucket]
    not matchBucketVersion(s3Bucket, bucketVersioning)

    bucketName := s3Bucket.config.bucket

    decode_rc := `resource "aws_s3_bucket_versioning" "example" {
    bucket = "aws_s3_bucket.##name##.bucket"

    versioning_configuration {
        status = "Enabled"
    }
    }`

    replacedName := replace(decode_rc, "##name##", bucketName)

    resourceWithoutQuotes := replace(replacedName, "\"", " ")

    traverse = ""
    retVal := getRetVal(s3Bucket.id, "add", "resource", traverse, "", "base64", base64.encode(resourceWithoutQuotes), null)
}

#To satisfy specific s3 buckets which have versioning configuration resource with status not enabled (for tfplan)
{{.prefix}}s3Versioning[retVal] {
    buck := input.aws_s3_bucket[_]
    object.get(buck.config, "versioning", "undefined") == [[], null, "undefined"][_]
    object.get(input, "aws_s3_bucket_versioning", "undefined") != "undefined"
   	bucketVersion := input.aws_s3_bucket_versioning[_]
    not contains(bucketVersion, "{")
    bucketVersioning :=  bucketVersion.config.bucket
    matchBucketVersion(buck, bucketVersioning)

    some i
    lower(bucketVersion.config.versioning_configuration[i].status) != "enabled"
    traverse := sprintf("versioning_configuration[%d].status", [i])

    retVal := getRetVal(bucketVersion.id, "edit", "attribute", traverse, "versioning_configuration.status", "string", "Enabled", bucketVersion.config.versioning_configuration[i].status )
}

#To satisfy specific s3 buckets which have versioning configuration resource with status not enabled (without tfplan)
{{.prefix}}s3Versioning[retVal] {
    buck := input.aws_s3_bucket[_]
    object.get(buck.config, "versioning", "undefined") == [[], null, "undefined"][_]
    object.get(input, "aws_s3_bucket_versioning", "undefined") != "undefined"
    bucketVersion := input.aws_s3_bucket_versioning[_]
    not contains(bucketVersion, "{")
    bucketVersioning := cleanSSEBucketID(bucketVersion.config.bucket)
    matchBucketVersion(buck, bucketVersioning)

    some i
    lower(bucketVersion.config.versioning_configuration[i].status) != "enabled"
    traverse := sprintf("versioning_configuration[%d].status", [i])

    retVal := getRetVal(bucketVersion.id, "edit", "attribute", traverse, "versioning_configuration.status", "string", "Enabled", bucketVersion.config.versioning_configuration[i].status )
}

matchBucketVersion(buckt, bucketVer){
   name := bucketVer[_]
   contains(name, buckt.name)
}

matchBucketVersion(buckt, bucketVer){
    buckt.config.bucket == bucketVer[_]
}

matchBucketVersion(buckt, bucketVer){
    buckt.config.bucket == bucketVer
}

matchBucketVersion(buckt, bucketVer){
    buckt.config.id == bucketVer[_]
}

matchBucketVersion(buckt, bucketVer){
    buckt.config.id == bucketVer
}

matchBucketVersion(buckt, bucketVer){
    buckt.id == bucketVer[_]
}

matchBucketVersion(buckt, bucketVer){
    buckt.id == bucketVer
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
        "Actual": Actual
    }
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
