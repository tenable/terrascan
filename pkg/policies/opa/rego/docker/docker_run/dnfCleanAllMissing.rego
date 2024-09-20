package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]
{
    run := input.docker_RUN[_]
    config := run.config
    
    installCommandExists(config)
    not contains(config, "dnf clean")
    not cleanExistsPostInstall(config)
}

installCommandExists(command) {
	installCommands = [
		"dnf install",
		"dnf in",
		"dnf reinstall",
		"dnf rei",
		"dnf install-n",
		"dnf install-na",
		"dnf install-nevra",
	]

	contains(command, installCommands[_])
}

cleanExistsPostInstall(config) {
	contains(config, "dnf clean all")

	installCommands = [
		"dnf install",
		"dnf in",
		"dnf reinstall",
		"dnf rei",
		"dnf install-n",
		"dnf install-na",
		"dnf install-nevra",
	]

	some cmd
	install := indexof(config, installCommands[cmd])
	install != -1

	clean := indexof(config, "dnf clean")
	clean != -1

	install < clean
}