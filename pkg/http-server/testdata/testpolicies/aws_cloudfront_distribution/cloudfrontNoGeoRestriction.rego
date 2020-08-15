package accurics

{{.prefix}}cloudfrontNoGeoRestriction[retVal] {
    cloudfront = input.aws_cloudfront_distribution[_]
    some i
    restrict = cloudfront.config.restrictions[i]
    restrict.geo_restriction[j].restriction_type == "none"
	traverse := sprintf("restrictions[%d].geo_restriction[%d].restriction_type", [i])
    retVal := { "Id": cloudfront.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "restrictions.geo_restriction.restriction_type", "AttributeDataType": "string", "Expected": "whitelist", "Actual": restrict.geo_restriction[_].restriction_type }
}