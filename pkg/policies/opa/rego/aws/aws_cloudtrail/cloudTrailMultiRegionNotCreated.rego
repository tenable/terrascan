package accurics

{{.prefix}}cloudTrailMultiRegionNotCreated[retVal]{
    cloud_trail = input.aws_cloudtrail[_]
    cloud_trail.config.is_multi_region_trail == false

    traverse = "is_multi_region_trail"
    retVal := { "Id": cloud_trail.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "is_multi_region_trail", "AttributeDataType": "bool", "Expected": true, "Actual": cloud_trail.config.is_multi_region_trail }
}