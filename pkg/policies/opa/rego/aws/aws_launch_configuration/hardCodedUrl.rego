package accurics

{{.prefix}}hardCodedUrl[res.id]{
    res = input.aws_instance[_]
    value = base64NullCheck(res.config.user_data_base64)
    contains(value, "https://")
}

{{.prefix}}hardCodedUrl[res.id]{
    res = input.aws_launch_configuration[_]
    value = base64NullCheck(res.config.user_data_base64)
    contains(value, "http://")
}

base64NullCheck(s) = result {
	s == null
	result := base64.decode("e30=")
}

base64NullCheck(s) = result {
	s != null
	result := base64.decode(s)
}