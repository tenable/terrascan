package accurics

{{.prefix}}{{.name}}{{.suffix}}[apt.config]{
	apt := input.run[_]
	#conval := apt.config
    containsCommand(apt.config)    
}    

hasInstall(cmds) {
	is_array(cmds)
	contains(cmds[_], "install")
}

hasInstall(cmd) {
	is_string(cmd) 
	contains(cmd, "install")
}

containsCommand(cmds) {
	count(cmds) > 1
	not hasInstall(cmds)
	regex.match("\\b(ps|shutdown|service|free|top|kill|mount|ifconfig|nano|vim)\\b", cmds[_])
}

containsCommand(cmds) {
	count(cmds) == 1

	commandsList = split(cmds[0], "&&")

	some i
	not hasInstall(commandsList[i])
	regex.match("\\b(ps|shutdown|service|free|top|kill|mount|ifconfig|nano|vim)\\b ", commandsList[i])
}

containsCommand(cmds) {
	count(cmds) == 1

	commandsList = split(cmds[0], "&&")

	some i
	not hasInstall(commandsList[i])
	regex.match("^\\b(ps|shutdown|service|free|top|kill|mount|ifconfig|nano|vim)\\b$", commandsList[i])
}

containsCommand(cmd) {
	string_encoded := json.marshal(cmd)
    not regex.match("(^[{[]).*([]}]$)", string_encoded)
	commandsList = split(cmd, " ")
    some i
    not hasInstall(commandsList[i])
	regex.match("^\\b(ps|shutdown|service|free|top|kill|mount|ifconfig|nano|vim)\\b$", commandsList[i])
}
