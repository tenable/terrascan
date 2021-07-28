package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.run[_]
    config := cmd.config
    contains(config, "apt")
}