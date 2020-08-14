package accurics

{{.prefix}}rdsHostsHigherThan256[retVal] {
    sg = input.aws_db_security_group[_]
    some i
    ingress = sg.config.ingress[i]
    hosts = split(ingress.cidr, "/")
	to_number(hosts[1]) <= 24
    expected := sprintf("%s/%d", [hosts[0], 23])
    traverse := sprintf("ingress[%d].cidr", [i])
    retVal := { "Id": sg.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ingress.cidr", "AttributeDataType": "list", "Expected": expected, "Actual": ingress.cidr }
}