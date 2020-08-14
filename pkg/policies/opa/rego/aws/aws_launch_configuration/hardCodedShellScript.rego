package accurics

{{.prefix}}hardCodedShellScript[res.id]{
    res = input.aws_instance[_]
    value = base64NullCheck(res.config.user_data_base64)
    startswith(value, "#!/")
}

{{.prefix}}hardCodedShellScript[res.id]{
    res = input.aws_launch_configuration[_]
    value = base64NullCheck(res.config.user_data_base64)
    startswith(value, "#!/")
}

base64NullCheck(s) = result {
    s == null
    result := base64.decode("e30=")
}

base64NullCheck(s) = result {
    s != null
    result := base64.decode(s)
}