package accurics

{{.prefix}}{{.name}}{{.suffix}}[retval] {
    sg := input.aws_security_group[_]
    some i
    ingress := sg.config.ingress[i]

    expected := checkConfig(ingress)
    traverse := sprintf("ingress[%d].cidr_blocks", [i])
    attribute := "ingress.cidr_blocks"

    retval := getretval(sg.id, traverse, attribute, expected, ingress.cidr_blocks)
}

{{.prefix}}{{.name}}{{.suffix}}[retval] {
    sgr := input.aws_security_group_rule[_]

    sgr.config.type == "ingress"
    expected := checkConfig(sgr.config)
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

checkConfig(config) = expected {
    config.cidr_blocks[_] == "0.0.0.0/0"
    checkProtocol(config.protocol, "{{.protocol}}")
    checkPort(config, {{.portNumber}})
    expected := [ item | item := validate_cidr(config.cidr_blocks[_]) ]
    expected != []
}

checkProtocol(configProtocol, protocol) {
    protocols = [protocol, "-1"]
    upper(configProtocol) == upper(protocols[_])
}

checkPort(config, port) {
    config.from_port == port
}

checkPort(config, port) {
    config.to_port == port
}

validate_cidr(cidr) = "{{.defaultValue}}" {
    cidr == "0.0.0.0/0"
}

validate_cidr(cidr) = cidr {
    cidr != "0.0.0.0/0"
}