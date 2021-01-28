package accurics

securityGroupWideOpentoInternetLaunchConfig[res.id] {
    res = input.aws_launch_configuration[_]
    sec_groups := [ sg | sg := split(res.config.security_groups[_], ".")[1] ]
    sec_group := sec_groups[_]
    checkSecurityGroupWideOpen(sec_group)
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