package accurics

{{.prefix}}rdsIsPublic[retVal] {
    sg = input.aws_db_security_group[_]
    some i
    ingress = sg.config.ingress[i]
    ingress.cidr == "0.0.0.0/0"
    traverse := sprintf("ingress[%d].cidr", [i])
    retVal := { "Id": sg.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ingress.cidr", "AttributeDataType": "list", "Expected": "<ingress.cidr>", "Actual": ingress.cidr }
}