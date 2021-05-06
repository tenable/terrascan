package accurics

{{.prefix}}cloudWatchLogGroupNoRetentionPolicy[cloudwatch.id] {
    cloudwatch := input.aws_cloudwatch_log_group[_]
    object.get(cloudwatch.config, "retention_in_days", "undefined") == "undefined"
}