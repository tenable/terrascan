package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]
{
run := input.run[_]
config := run.config
isZypperUnsafeCommand(config)
}

isZypperUnsafeCommand(command) {
	contains(command, "zypper update")
}

isZypperUnsafeCommand(command) {
	contains(command, "zypper dist-upgrade")
}

isZypperUnsafeCommand(command) {
	contains(command, "zypper dup")
}

isZypperUnsafeCommand(command) {
	contains(command, "zypper up")
}
