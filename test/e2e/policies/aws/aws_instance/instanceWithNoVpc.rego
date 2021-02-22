package accurics

{{.prefix}}instanceWithNoVpc[retVal] {
    instance := input.aws_instance[_]
    not instance.config.vpc_security_group_ids
    rc = "ewogICJhd3NfdnBjIjogewogICAgImFjY3VyaWNzX3ZwYyI6IHsKICAgICAgImNpZHJfYmxvY2siOiAiPGNpZHJfYmxvY2s+IiwKICAgICAgImVuYWJsZV9kbnNfc3VwcG9ydCI6ICI8ZW5hYmxlX2Ruc19zdXBwb3J0PiIsCiAgICAgICJlbmFibGVfZG5zX2hvc3RuYW1lcyI6ICI8ZW5hYmxlX2Ruc19ob3N0bmFtZXM+IgogICAgfQogIH0KfQ=="
    traverse = ""
    retVal := { "Id": instance.id, "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "resource", "Expected": rc, "Actual": null }
}