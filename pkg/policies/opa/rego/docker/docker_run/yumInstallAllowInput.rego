package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
	cmd := input.docker_RUN[_]
    config := cmd.config
    checkYumInstall(config)
    not checkManualInput(config)
}

checkYumInstall(config) {
	re_match(`yum (-(-)?[a-zA-Z]+ *)*(group|local)?install`, config)
}

checkManualInput(config) {
	commands := ["-y", "yes", "--assumeyes"]
    checkCmd := commands[_]
    contains(config, checkCmd)
}