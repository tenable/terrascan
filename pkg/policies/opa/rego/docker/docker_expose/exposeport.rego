package accurics

{{.prefix}}{{.name}}{{.suffix}}[expose.id]{
	expose := input.docker_expose[_]
	is_string(expose.config)
    config := expose.config
    checkPort(config)
}

{{.prefix}}{{.name}}{{.suffix}}[expose.id] {
    expose := input.docker_expose[_]
    is_array(expose.config)
    config := expose.config
    checkPortList(config)
}

checkPort(config) {
    contains(config, "22")
}

checkPortList(config) {
    contains(config[_], "22")
}  