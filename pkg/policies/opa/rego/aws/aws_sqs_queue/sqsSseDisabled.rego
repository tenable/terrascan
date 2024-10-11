package accurics

{{.prefix}}sqsSseDisabled[sqs.id] {
  sqs := input.aws_sqs_queue[_]
  object.get(sqs.config, "kms_master_key_id", "undefined") == ["undefined", null, ""][_]
  object.get(sqs.config, "sqs_managed_sse_enabled", false) == [false , null, ""][_]
}
