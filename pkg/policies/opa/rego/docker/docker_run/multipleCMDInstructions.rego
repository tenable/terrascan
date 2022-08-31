package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmdInst[_]]{
	command := input.docker_cmd
    cmdInst := [x | x := checkCommandType(command[_])]
	count(cmdInst) > 1
}

checkCommandType(command) = value {
	command.type == "docker_cmd"
    value := command.id
}