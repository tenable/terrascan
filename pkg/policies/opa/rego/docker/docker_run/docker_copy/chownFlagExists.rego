package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.copy[_]
    config := cmd.config
    contains(config, "--chown")
}