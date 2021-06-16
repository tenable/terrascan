package accurics

{{.prefix}}{{.name}}{{.suffix}}[retval] {
    security_group := input.aws_security_group[_]

    some i
    some j
    ingress := security_group.config.ingress[i]

    ingress.cidr_blocks[j] == "0.0.0.0/0"
    check_config(ingress)

    expected := [ item | item := validate_cidr(ingress.cidr_blocks[_]) ]
    traverse := sprintf("ingress[%d].cidr_blocks", [i])
    attribute := "ingress.cidr_blocks"

    retval := getretval(security_group.id, traverse, attribute, expected, ingress.cidr_blocks)
}

{{.prefix}}{{.name}}{{.suffix}}[retval] {
    sgr := input.aws_security_group_rule[_]
    sgr.config.type == "ingress"

    some i
    cidr := sgr.config.cidr_blocks[i]

    cidr == "0.0.0.0/0"
    check_config(sgr.config)

    expected := [ item | item := validate_cidr(sgr.config.cidr_blocks[_]) ]
    traverse_attribute := "cidr_blocks"

    retval := getretval(sgr.id, traverse_attribute, traverse_attribute, expected, sgr.config.cidr_blocks)
}

getretval(id, traverse, attribute, expected, actual) = retval {
    retval := {
        "Id": id,
        "ReplaceType": "edit",
        "CodeType": "attribute",
        "Traverse": traverse,
        "Attribute": attribute,
        "AttributeDataType": "list",
        "Expected": expected,
        "Actual": actual
    }
}

check_config(config) {
    config.from_port == 0
    config.to_port == 0
    config.protocol == "-1"
}

validate_cidr(cidr) = "{{.defaultValue}}" {
    cidr == "0.0.0.0/0"
}

validate_cidr(cidr) = cidr {
    cidr != "0.0.0.0/0"
}