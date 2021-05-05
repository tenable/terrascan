package accurics

{{.prefix}}noSgDescription[sg.id]{
    sg = input.aws_security_group[_]
	object.get(sg.config, "description", "undefined") = ["undefined", "Managed by Terraform"][_]
}