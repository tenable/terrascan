package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.docker_RUN[_]
    config := cmd.config
    contains(config, "apt")
}