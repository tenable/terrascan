package accurics

{{.prefix}}cloudTrailLogValidationDisabled[cloudtrail.id] {
    cloudtrail := input.aws_cloudtrail[_]
    object.get(cloudtrail.config, "enable_log_file_validation", "undefined") == "undefined"
}

{{.prefix}}cloudTrailLogValidationDisabled[cloudtrail.id] {
    cloudtrail := input.aws_cloudtrail[_]
    cloudtrail.config.enable_log_file_validation == false
}