package accurics

{{.prefix}}kinesisNotEncryptedWithKms[stream.id] {
    stream = input.aws_kinesis_stream[_]
    stream.config.kms_key_id == null
}

{{.prefix}}kinesisNotEncryptedWithKms[stream.id] {
    stream = input.aws_kinesis_stream[_]
    object.get(stream.config, "encryption_type", "undefined") == ["NONE", "undefined"][_]
}