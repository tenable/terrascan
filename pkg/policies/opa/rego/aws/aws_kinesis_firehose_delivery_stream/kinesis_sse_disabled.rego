package accurics

{{.prefix}}kinesisSseDisabled[retVal] {
    stream = input.aws_kinesis_firehose_delivery_stream[_]
    some i
	stream.config.server_side_encryption[i].enabled == false
    traverse := sprintf("server_side_encryption[%d].enabled", [i])
    retVal := { "Id": stream.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "server_side_encryption.enabled", "AttributeDataType": "bool", "Expected": true, "Actual": stream.config.server_side_encryption[i].enabled }
}
