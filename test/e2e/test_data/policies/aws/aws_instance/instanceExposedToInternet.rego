package accurics

instanceExposedToInternet[ins.id] {
    ins = input.aws_instance[_]
    sec_groups := [ sg | sg := split(ins.config.vpc_security_group_ids[_], ".")[1] ]
    sec_group := sec_groups[_]
    checkSecurityGroupWideOpen(sec_group)

    sub = ins.config.subnet_id
    route_table_association = input.aws_route_table_association[_]
    route_table_association.config.subnet_id == sub
    route_table := split(route_table_association.config.route_table_id, ".")[1]
    checkRouteInternet(ins, route_table)
}

checkSecurityGroupWideOpen(sgName) {
	security_group := input.aws_security_group[_]
    sgName == security_group.name

    some i
    ingress = security_group.config.ingress[i]

    # Checks if the cidr block is not a private IP
    ingress.cidr_blocks[_] == "0.0.0.0/0"

    ports_open = (ingress.to_port - ingress.from_port)
    ports_open > 0
}

checkRouteInternet(instance, arg) {
    rt = input.aws_route_table[_]
    rt.name == arg
    routes = rt.config.route[_]
    routes.gateway_id != ""
}