package accurics

{{.prefix}}{{.name}}{{.suffix}}[nsg.id] {
	nsg := input.azurerm_network_security_group[_]

	object.get(nsg.config, "security_rule", "undefined") != "undefined"

	some i
	sg := nsg.config.security_rule[i]
	sg.access == "Allow"
	sg.direction == "Inbound"

	source_prefix := sg.source_address_prefixes[_]
    not re_match(`[a-zA-Z]+`, source_prefix)

	private_ips = ["10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"]
	net.cidr_contains(private_ips[_], sg.source_address_prefixes[_])
	count(sg.source_address_prefixes) > 32
}

{{.prefix}}{{.name}}{{.suffix}}[nsg.id] {
	nsg := input.azurerm_network_security_group[_]

	object.get(nsg.config, "security_rule", "undefined") != "undefined"

	some i
	sg := nsg.config.security_rule[i]
	sg.access == "Allow"
	sg.direction == "Inbound"

	source_prefix := sg.source_address_prefixes[_]
    not re_match(`[a-zA-Z]+`, source_prefix)

	private_ips = ["10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"]
	net.cidr_contains(private_ips[_], source_prefix)
	hosts = split(source_prefix, "/")
    to_number(hosts[1]) < 27
}
