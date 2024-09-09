package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.docker_RUN[_]
    config := cmd.config
    commands = ["dnf update", "dnf upgrade", "dnf upgrade-minimal"]
    contains(config, commands[_])
}