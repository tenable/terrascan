package accurics

{{.prefix}}kinesisNotEncryptedWithKms[retVal] {
    stream = input.aws_kinesis_stream[_]
    stream.config.kms_key_id == null
    traverse = ""
    retVal := { "Id": stream.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "kms_key_id", "AttributeDataType": "string", "Expected": "<kms_key_id>", "Actual": null }
}