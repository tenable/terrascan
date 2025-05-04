package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.docker_FROM[_]
    config := cmd.config
    contains(config, "--platform")
}