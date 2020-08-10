package accurics

{{.prefix}}enableSNSTopic[sns.id] {
    sns := input.aws_cloudtrail[_]
    sns.config.sns_topic_name == null
}