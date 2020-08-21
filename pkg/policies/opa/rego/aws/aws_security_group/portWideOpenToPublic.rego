package accurics

{{.prefix}}portWideOpenToPublic[retVal] {
    sg = input.aws_security_group[_]
    some i
    ingress = sg.config.ingress[i]

    # Checks if the cidr block is not a private IP
    ingress.cidr_blocks[_] == "0.0.0.0/0"

    ports_open = (ingress.to_port - ingress.from_port)

    ports_open > 0

    traverse := sprintf("ingress[%d].to_port", [i])
    expected := ingress.to_port

    retVal := { "Id": sg.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "ingress.from_port", "AttributeDataType": "int", "Expected": expected, "Actual": ingress.from_port }
}