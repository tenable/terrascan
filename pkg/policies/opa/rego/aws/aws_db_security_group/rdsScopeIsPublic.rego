package accurics

{{.prefix}}rdsScopeIsPublic[retVal] {
    sg = input.aws_db_security_group[_]
    some i
    ingress = sg.config.ingress[i]
    checkScopeIsPublic(ingress.cidr)
    traverse := sprintf("ingress[%d].cidr", [i])
	retVal := { "Id": sg.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ingress.cidr", "AttributeDataType": "list", "Expected": "<ingress.cidr>", "Actual": ingress.cidr }
}

scopeIsPrivate(scope) {
    private_ips = ["10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"]
    net.cidr_contains(private_ips[_], scope)
}

checkScopeIsPublic(val) {
    not scopeIsPrivate(val)
}