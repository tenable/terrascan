package accurics

{{.prefix}}{{.name}}[retVal] {
  sg = input.azurerm_network_security_rule[_]
  sg.config.access == "Allow"
  sg.config.direction == "Inbound"
  checkScopeIsPublic(sg.config.source_address_prefix)
  checkPort(sg.config, "{{.portNumber}}")
  checkProtocol(sg.config.protocol)

  traverse := "source_address_prefix"
  retVal := { "Id": sg.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "source_address_prefix", "AttributeDataType": "string", "Expected": "<cidr>", "Actual": sg.config.source_address_prefix }
}

{{.prefix}}{{.name}}[retVal] {
  nsg = input.azurerm_network_security_group[_]
  some i
  sg = nsg.config.security_rule[i]
  sg.access == "Allow"
  sg.direction == "Inbound"
  checkScopeIsPublic(sg.source_address_prefix)
  checkPort(sg, "{{.portNumber}}")
  checkProtocol(sg.protocol)

  traverse := sprintf("security_rule[%d].source_address_prefix", [i])
  retVal := { "Id": nsg.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "security_rule.source_address_prefix", "AttributeDataType": "string", "Expected": "<cidr>", "Actual": sg.source_address_prefix }
}

scopeIsPrivate(scope) {
  private_ips = ["10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"]
  net.cidr_contains(private_ips[_], scope)
}

checkScopeIsPublic(val) = true {
  glob.match("[0-9]*.[0-9]*.[0-9]*.*", [], val)
  not scopeIsPrivate(val)
  hosts = split(val, "/")
  to_number(hosts[1]) < {{.numberOfHosts}}
  to_number(hosts[1]) >= {{.endLimit}}
}

checkScopeIsPublic(val) = true {
  glob.match("[0-9]*.[0-9]*.[0-9]*.*", [], val)
  not scopeIsPrivate(val)
  hosts = split(val, "/")
  not hosts[1]
  {{.evalHosts}}
}

checkScopeIsPublic(val) = true {
  not glob.match("[0-9]*.[0-9]*.[0-9]*.*", [], val)
  val == "*"
  {{.evalHosts}}
}

checkScopeIsPublic(val) = true {
  not glob.match("[0-9]*.[0-9]*.[0-9]*.*", [], val)
  val == "Internet"
  {{.evalHosts}}
}

checkPort(obj, val) = true {
  obj.destination_port_range == val
}

checkPort(obj, val) = true {
  obj.source_port_range == val
}

checkProtocol(proto) {
  protocols = ["{{.protocol}}", "*"]
  upper(proto) == protocols[_]
}
