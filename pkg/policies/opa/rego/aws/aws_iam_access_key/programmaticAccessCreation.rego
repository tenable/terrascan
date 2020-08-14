package accurics

{{.prefix}}programmaticAccessCreation[retVal] {
    access := input.aws_iam_access_key[_]
    access.type == "aws_iam_access_key"
	status = getStatus(access.config)
    status == "Active"
    access.config.user != "root"
    traverse = "status"
    retVal := { "Id": access.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "status", "AttributeDataType": "string", "Expected": "Inactive", "Actual": status }

}

getStatus(config) = "Active" {
    # defaults to Active
    not config.status
}

getStatus(config) = "Active" {
    config.status == "Active"
}