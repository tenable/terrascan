package accurics

{{.prefix}}noRuleDescription[sg.id]{
    sg = input.aws_security_group[_]
	egress := sg.config.egress[_]
    egress.description == ["", " "][_] #for terraformer quotes have a space
}

{{.prefix}}noRuleDescription[sg.id]{
    sg = input.aws_security_group[_]
	ingress := sg.config.ingress[_]
    ingress.description == ["", " "][_]
}