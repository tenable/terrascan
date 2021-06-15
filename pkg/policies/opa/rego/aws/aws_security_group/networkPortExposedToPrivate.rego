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
        "Actual": actual,
    }
}

checkConfig(config) = expected {
    checkPort(config, {{.portNumber}})
    checkProtocol(config.protocol, "{{.protocol}}")
    expected := [item | item := checkScopeIsPrivate(config.cidr_blocks[_])]
    expected != []
}

checkPort(config, port) {
    config.from_port == port
}

checkPort(config, port) {
    config.to_port == port
}

checkProtocol(configProtocol, protocol) {
    protocols = [protocol, "-1"]
    upper(configProtocol) == upper(protocols[_])
}

checkScopeIsPrivate(ingress_cidr) = value {
    glob.match("[0-9]*.[0-9]*.[0-9]*.*", [], ingress_cidr)

    private_ips = ["10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"]
    net.cidr_contains(private_ips[_], ingress_cidr)

    hosts = split(ingress_cidr, "/")
    to_number(hosts[1]) < 27

    value := "{{.defaultValue}}"
}