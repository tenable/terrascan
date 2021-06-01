package accurics

{{.prefix}}cloudTrailMultiRegionEnabled[cloud_trail.id]{
    cloud_trail = input.aws_cloudtrail[_]
    object.get(cloud_trail, "is_multi_region_trail", "undefined") == "undefined"
}