package accurics

{{.prefix}}cloudTrailLogNotEncrypted[cloud_trail.id]{
    cloud_trail = input.aws_cloudtrail[_]
    object.get(cloud_trail.config, "kms_key_id", "undefined") == [null, "undefined"][_]
}