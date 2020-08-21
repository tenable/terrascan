package accurics

{{.prefix}}{{.name}}[retVal]{
    security_group = input.aws_security_group[_]
    some i
    ingress = security_group.config.ingress[i]
    ingress.protocol == "{{.protocol}}"
    ingress.from_port == {{.from_port}}
    ingress.cidr_blocks[j] == "0.0.0.0/0"
    expected := [ item | item := validate_cidr(ingress.cidr_blocks[_]) ]
    traverse := sprintf("ingress[%d].cidr_blocks", [i])
    retVal := { "Id": security_group.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ingress.cidr_blocks", "AttributeDataType": "list", "Expected": expected, "Actual": ingress.cidr_blocks }
}

validate_cidr(cidr) = value {
	cidr == "0.0.0.0/0"
    value := "<cidr>"
}

validate_cidr(cidr) = value {
	cidr != "0.0.0.0/0"
    value := cidr
}