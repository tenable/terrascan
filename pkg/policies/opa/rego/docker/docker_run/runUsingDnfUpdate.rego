package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.run[_]
    config := cmd.config
    commands = ["dnf update", "dnf upgrade", "dnf upgrade-minimal"]
    contains(config, commands[_])
}