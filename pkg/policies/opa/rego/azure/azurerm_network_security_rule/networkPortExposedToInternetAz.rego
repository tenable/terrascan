package accurics

{{.prefix}}{{.name}}{{.suffix}}[retval] {
    sg := input.azurerm_network_security_rule[_]
    checkConfiguration(sg.config)

    traverse_attribute := "source_address_prefix"
    retval := getretval(sg.id, traverse_attribute, traverse_attribute, sg.config.source_address_prefix)
}

{{.prefix}}{{.name}}{{.suffix}}[retval] {
    nsg := input.azurerm_network_security_group[_]

    object.get(nsg.config, "security_rule", "undefined") != "undefined"

    some i
    sg := nsg.config.security_rule[i]
    checkConfiguration(sg)

    traverse := sprintf("security_rule[%d].source_address_prefix", [i])
    attribute := "security_rule.source_address_prefix"
    retval := getretval(nsg.id, traverse, attribute, sg.source_address_prefix)
}

getretval(id, traverse, attribute, actual) = retval {
    retval := {
        "Id": id,
        "ReplaceType": "edit",
        "CodeType": "attribute",
        "Traverse": traverse,
        "Attribute": attribute,
        "AttributeDataType": "string",
        "Expected": "{{.defaultValue}}",
        "Actual": actual
    }
}

checkConfiguration(sg) {
    sg.access == "Allow"
    sg.direction == "Inbound"

    checkOpenToInternet(sg.source_address_prefix)
    checkPort(sg, "{{.portNumber}}")
    checkProtocol(sg.protocol, "{{.protocol}}")
}

checkOpenToInternet(source_address_prefix) {
    source_address_prefix == "*"
}

checkOpenToInternet(source_address_prefix) {
    source_address_prefix == "0.0.0.0/0"
}

checkPort(config, port) {
    config.destination_port_range == port
}

checkPort(config, port) {
    config.source_port_range == port
}

checkProtocol(configProtocol, protocol) {
    protocols = [protocol, "*"]
    upper(configProtocol) == upper(protocols[_])
}
