package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.docker_run[_]
    config := cmd.config
    contains(config, "apt")
}