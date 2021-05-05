package accurics

{{.prefix}}notEncryptedSns[sns_topic.id] {
    sns_topic := input.aws_sns_topic[_]
    object.get(sns_topic.config, "kms_master_key_id", "undefined") == [null, "undefined"][_]
}