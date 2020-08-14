package accurics

{{.prefix}}hardCodedKey[res.id] {
    res = input.aws_launch_configuration[_]
	value = base64NullCheck(res.config.user_data_base64)
	contains(value, "LS0tLS1CR")
}

{{.prefix}}hardCodeKey[res.id]{
    res = input.aws_launch_configuration[_]
    value = base64NullCheck(res.config.user_data_base64)
    contains(value, "LS0tLS1CR")
}

base64NullCheck(s) = result {
	s == null
	result := base64.decode("e30=")
}

base64NullCheck(s) = result {
	s != null
	result := base64.decode(s)
}