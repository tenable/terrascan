package accurics

{{.prefix}}cloudfrontNoHTTPSTraffic[retVal]{
    cloudfront = input.aws_cloudfront_distribution[_]
    some i
    orderedcachebehaviour = cloudfront.config.ordered_cache_behavior[i]
    orderedcachebehaviour.viewer_protocol_policy == "allow-all"
	traverse := sprintf("ordered_cache_behavior[%d].viewer_protocol_policy", [i])
    retVal := { "Id": cloudfront.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ordered_cache_behavior.viewer_protocol_policy", "AttributeDataType": "string", "Expected": "redirect-to-https", "Actual": orderedcachebehaviour.viewer_protocol_policy }
}