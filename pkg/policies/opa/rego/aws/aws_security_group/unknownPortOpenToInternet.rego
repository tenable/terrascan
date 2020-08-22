package accurics

{{.prefix}}unknownPortOpenToInternet[retVal]{
    security_group = input.aws_security_group[_]
    some i
    ingress = security_group.config.ingress[i]

    ingress.cidr_blocks[j] == "0.0.0.0/0"
    known_ports = [{{range .known_ports}}{{- printf "%q" . }},{{end}}]
    not contains_port(known_ports, ingress.from_port)

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

contains_port(known_ports, port) {
	known_ports[_] == sprintf("%d", [port])
}