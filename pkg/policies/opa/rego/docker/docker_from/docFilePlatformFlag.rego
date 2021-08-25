package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.docker_from[_]
    config := cmd.config
    contains(config, "--platform")
}