package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]
{
    run := input.docker_run[_]
    config := run.config
    command := config
    
    validUpdateCommand(command)
	not updateAfterInstall(command)
}

validUpdateCommand(command) {
	contains(command, " update ")
}

validUpdateCommand(command) {
	contains(command, " --update ")
}

validUpdateCommand(command) {
	array_split := split(command, " ")

	len = count(array_split)

	update := {"update", "--update"}

	array_split[minus(len, 1)] == update[j]
}

updateAfterInstall(command) {
	commandList = [
		"install",
		"source-install",
		"reinstall",
		"groupinstall",
		"localinstall",
		"add",
	]

	update := indexof(command, "update")
	update != -1

	install := indexof(command, commandList[_])
	install != -1

	update < install
}