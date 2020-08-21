package accurics

{{.prefix}}defaultSGNotRestrictsAllTraffic[retVal] {
    security_group = input.aws_security_group[_]

    disabled = true

    not security_group.config.ingress
    not security_group.config.egress

    # This Rule will conflict security_group_rule and must be disabled
    # This is not a valid rule anyways, because if you create a security group
    # without ingress or egress, AWS would not allow you to select it.
    disabled == false

    rc = "ewogICJpbmdyZXNzIjogewogICAgInRvX3BvcnQiOiAiPGluZ3Jlc3NfdG9fcG9ydD4iLAogICAgImZyb21fcG9ydCI6ICI8aW5ncmVzc19mcm9tX3BvcnQ+IiwKICAgICJwcm90b2NvbCI6ICI8aW5ncmVzc19wcm90b2NvbD4iLAogICAgImNpZHJfYmxvY2tzIjogWyI8aW5ncmVzc19jaWRyX2Jsb2Nrcz4iXQogIH0sCiAgImVncmVzcyI6IHsKICAgICJmcm9tX3BvcnQiOiAiPGVncmVzc19mcm9tX3BvcnQ+IiwKICAgICJ0b19wb3J0IjogIjxlZ3Jlc3NfdG9fcG9ydD4iLAogICAgInByb3RvY29sIjogIjxlZ3Jlc3NfcHJvdG9jb2w+IiwKICAgICJjaWRyX2Jsb2NrcyI6IFsiPGVncmVzc19jaWRyX2Jsb2Nrcz4iXQogIH0KfQ=="
    traverse = ""
    retVal := { "Id": security_group.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "", "AttributeDataType": "block", "Expected": rc, "Actual": null }
}