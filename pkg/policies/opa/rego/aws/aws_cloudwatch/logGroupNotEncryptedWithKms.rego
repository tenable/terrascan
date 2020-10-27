package accurics

logGroupNotEncryptedWithKms[log_group.id] {
  log_group := input.aws_cloudwatch_log_group[_]
  not log_group.config.kms_key_id
}
