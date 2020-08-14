package accurics

{{.prefix}}kinesisSseNotConfigured[retVal] {
    stream := input.aws_kinesis_firehose_delivery_stream[_]
    count(stream.config.server_side_encryption) == 0
    rc = "ewogICJzZXJ2ZXJfc2lkZV9lbmNyeXB0aW9uIjogewogICAgImVuYWJsZWQiOiB0cnVlCiAgfQp9"
    traverse = ""
    retVal := { "Id": stream.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "server_side_encryption", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}
