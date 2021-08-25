package accurics

{{.prefix}}{{.name}}{{.suffix}}[apt.id]{
	apt := input.docker_expose[_]
	conval := apt.config
    port := split(conval, "/")
    containsPortOutOfRange(port)
}
containsPortOutOfRange(ports) {
	some i
	port := ports[i]
	to_number(port) > 65535
}