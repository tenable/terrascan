package accurics

{{.prefix}}{{.name}}[retVal] {
    sg = input.aws_security_group[_]
    some i
    ingress = sg.config.ingress[i]
    # Checks if the cidr block is not a private IP
    checkScopeIsPublic(ingress.cidr_blocks[j])
    checkProtocol(ingress.protocol)
    # Check if port range matches what we are detecting.
    checkPort(ingress, {{.portNumber}})

    expected := [ item | item := validate_cidr(ingress.cidr_blocks[_]) ]
    traverse := sprintf("ingress[%d].cidr_blocks", [i])
    retVal := { "Id": sg.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ingress.cidr_blocks", "AttributeDataType": "list", "Expected": expected, "Actual": ingress.cidr_blocks }
}

scopeIsPrivate(scope) {
    private_ips = ["10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"]
    net.cidr_contains(private_ips[_], scope)
}

checkScopeIsPublic(val) {
    not scopeIsPrivate(val)
    val != "0.0.0.0/0"
}

checkProtocol(proto) {
    protocols = ["{{.protocol}}", "-1"]
    proto == protocols[_]
}

checkPort(obj, val) = true {
    obj.from_port == val
    obj.to_port == val
}

validate_cidr(cidr) = value {
	checkScopeIsPublic(cidr)
    value := "<cidr>"
}

validate_cidr(cidr) = value {
	not checkScopeIsPublic(cidr)
    cidr == "0.0.0.0/0"
    value := "<cidr>"
}

validate_cidr(cidr) = value {
	not checkScopeIsPublic(cidr)
    cidr != "0.0.0.0/0"
    value := cidr
}