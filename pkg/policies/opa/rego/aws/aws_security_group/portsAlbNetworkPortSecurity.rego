package accurics

{{.prefix}}{{.name}}[retVal] {
    sg = input.aws_security_group[_]
    some i
    ingress = sg.config.ingress[i]
    # Checks if the cidr block is not a private IP
    ingress.cidr_blocks[j] == "0.0.0.0/0"
    checkProtocol(ingress.protocol)
    # Check if port range matches what we are detecting.
    checkPort(ingress, {{.portNumber}})

    expected := [ item | item := validate_cidr(ingress.cidr_blocks[_]) ]
    traverse := sprintf("ingress[%d].cidr_blocks", [i])
    retVal := { "Id": sg.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ingress.cidr_blocks", "AttributeDataType": "list", "Expected": expected, "Actual": ingress.cidr_blocks }
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
	cidr == "0.0.0.0/0"
    value := "<cidr>"
}

validate_cidr(cidr) = value {
	cidr != "0.0.0.0/0"
    value := cidr
}