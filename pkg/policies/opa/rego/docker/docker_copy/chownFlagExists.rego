package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.docker_COPY[_]
    config := cmd.config
    contains(config, "--chown")
}