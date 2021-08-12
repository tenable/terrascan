package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.docker_copy[_]
    config := cmd.config
    contains(config, "--chown")
}