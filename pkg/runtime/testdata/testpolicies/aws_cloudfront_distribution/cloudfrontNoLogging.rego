package accurics

{{.prefix}}cloudfrontNoLogging[retVal]{
    cloudfront = input.aws_cloudfront_distribution[_]
    not cloudfront.config.logging_config

    rc = "ewogICJsb2dnaW5nX2NvbmZpZyI6IHsKICAgICJpbmNsdWRlX2Nvb2tpZXMiOiBmYWxzZSwKICAgICJidWNrZXQiOiAiPGJ1Y2tldD4iLAogICAgInByZWZpeCI6ICI8cHJlZml4PiIKICB9Cn0="

    traverse = ""
    retVal := { "Id": cloudfront.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "logging_config", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}

{{.prefix}}cloudfrontNoLogging[retVal]{
    cloudfront = input.aws_cloudfront_distribution[_]
    cloudfront.config.logging_config == []

    rc = "ewogICJsb2dnaW5nX2NvbmZpZyI6IHsKICAgICJpbmNsdWRlX2Nvb2tpZXMiOiBmYWxzZSwKICAgICJidWNrZXQiOiAiPGJ1Y2tldD4iLAogICAgInByZWZpeCI6ICI8cHJlZml4PiIKICB9Cn0="

    traverse = ""
    retVal := { "Id": cloudfront.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "logging_config", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}