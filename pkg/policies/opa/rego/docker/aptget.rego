package accurics

{{.prefix}}{{.name}}{{.suffix}}[conval]
{
	apt := input.run[_]
	conval := apt.config
	commandsSplit := split(conval, "&&")
	not avoidAdditionalPackages(commandsSplit)
}
avoidAdditionalPackages(commands) {
	is_string(commands) == true
	flags := ["--no-install-recommends", "apt::install-recommends=false"]
    some i
    flag := flags[i]
	contains(commands, flag)
}
avoidAdditionalPackages(commands) {
   is_array(commands) == true
   some i
   command := commands[i]
   contains(command, "--no-install-recommends")
}
avoidAdditionalPackages(commands) {
   is_array(commands) == true
   some i
   command := commands[i]
   contains(command, "apt::install-recommends=false")  
}