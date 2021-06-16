package accurics

{{.prefix}}{{.name}}{{.suffix}}[retval] {
    sg = input.aws_security_group[_]
    some i
    ingress = sg.config.ingress[i]

    checkConfig(ingress)

    traverse := sprintf("ingress[%d].to_port", [i])
    attribute := "ingress.from_port"
    expected := ingress.to_port

    retval := getretval(sg.id, traverse, attribute, expected, ingress.from_port)
}

{{.prefix}}{{.name}}{{.suffix}}[retval] {
    sgr = input.aws_security_group_rule[_]
    sgr.config.type == "ingress"

    checkConfig(sgr.config)

    expected := sgr.config.to_port
    traverse_attribute = "from_port"

    retval := getretval(sgr.id, traverse_attribute, traverse_attribute, expected, sgr.config.from_port)
}

getretval(id, traverse, attribute, expected, actual) = retval {
    retval := {
        "Id": id,
        "ReplaceType": "edit",
        "CodeType": "attribute",
        "AttributeDataType": "int",
        "Traverse": traverse,
        "Attribute": attribute,
        "Expected": expected,
        "Actual": actual
    }
}

checkConfig(ingress) {
    ingress.cidr_blocks[_] == "0.0.0.0/0"

    ports_open = (ingress.to_port - ingress.from_port)
    ports_open > 0
}