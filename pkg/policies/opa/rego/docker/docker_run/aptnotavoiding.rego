package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerRun.name] {
    dockerRun := input.docker_RUN[_]
    is_string(dockerRun.config)
    config := dockerRun.config
    configArray := split(config, "&&")
    command := configArray[_]
    
    startswith(command, ["sudo apt-get", "apt-get"][_])
    contains(command, "install")

    not avoidAdditionalPackages(command)
}

avoidAdditionalPackages(arg) {
	flags := ["--no-install-recommends", "apt::install-recommends=false"]
	contains(arg, flags[_])
}