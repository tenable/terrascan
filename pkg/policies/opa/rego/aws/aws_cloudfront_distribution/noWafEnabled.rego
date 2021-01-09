package accurics

noWafEnabled[retVal] {
    cloudfront := input.aws_cloudfront_distribution[_]
    not cloudfront.config.web_acl_id

    traverse = "web_acl_id"
    retVal := { "Id": cloudfront.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "string", "Expected": "<arn-of-waf-acl>", "Actual": null }
}

noWafEnabled[retVal] {
    cloudfront := input.aws_cloudfront_distribution[_]
    cloudfront.config.web_acl_id == null

    traverse = "web_acl_id"
    retVal := { "Id": cloudfront.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "string", "Expected": "<arn-of-waf-acl>", "Actual": null }
}
