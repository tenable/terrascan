package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.docker_RUN[_]
    config := cmd.config
    regex.match(`(^|\W)apt\s`, config)
}
