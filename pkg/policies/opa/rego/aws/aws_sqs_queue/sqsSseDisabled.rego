package accurics

{{.prefix}}sqsSseDisabled[retVal] {
  sqs := input.aws_sqs_queue[_]
  check_empty(sqs.config.kms_master_key_id)
  traverse = "kms_master_key_id"
  retVal := { "Id": sqs.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "kms_master_key_id", "AttributeDataType": "string", "Expected": "<kms_master_key_id>" }
}

check_empty(key) {
	key == null
}
check_empty(key) {
	key == ""
}