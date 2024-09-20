package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmdInst[_]]{
	command := input.docker_CMD
    cmdInst := [x | x := checkCommandType(command[_])]
	count(cmdInst) > 1
}

checkCommandType(command) = value {
	command.type == "docker_CMD"
    value := command.id
}